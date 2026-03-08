use std::sync::Arc;

use async_trait::async_trait;
use axum::{
    extract::State,
    http::{header, StatusCode},
    middleware::{self, Next},
    response::{IntoResponse, Response},
    routing::post,
    Json, Router,
};
use axum_extra::extract::CookieJar;
use axum_extra::extract::cookie::{Cookie, SameSite};
use chrono::Utc;
use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};
use sea_orm::{
    ActiveModelTrait, ActiveValue::Set, ColumnTrait, DatabaseConnection, EntityTrait,
    QueryFilter,
};
use sea_orm_migration::MigratorTrait;
use serde::{Deserialize, Serialize};
use serde_json::json;

use backend_core::{crud_router, AppError, Module};

use crate::{
    entities::{
        role,
        user::{self, ActiveModel as UserActiveModel},
    },
    migration::Migrator,
    providers::{RoleProvider, UserProvider},
};

// ── JWT Claims ────────────────────────────────────────────────────────────────

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Claims {
    pub user_id: i64,
    pub email: String,
    pub role: String,
    pub exp: i64,
    pub iat: i64,
}

// ── Shared state ──────────────────────────────────────────────────────────────

#[derive(Clone)]
struct AuthState {
    db: DatabaseConnection,
    jwt_secret: String,
}

// ── AuthModule ────────────────────────────────────────────────────────────────

pub struct AuthModule {
    db: DatabaseConnection,
    jwt_secret: String,
}

impl AuthModule {
    pub fn new(db: DatabaseConnection, jwt_secret: String) -> Self {
        Self { db, jwt_secret }
    }

    /// Axum middleware that validates JWT from Authorization header or cookie,
    /// and injects Claims into request extensions.
    pub fn require_auth() -> axum::middleware::FromFnLayer<
        impl Fn(axum::http::Request<axum::body::Body>, Next) -> std::pin::Pin<Box<dyn std::future::Future<Output = Response> + Send>>,
        (),
        (axum::http::Request<axum::body::Body>, Next),
    > {
        middleware::from_fn(auth_middleware)
    }
}

#[async_trait]
impl Module for AuthModule {
    fn name(&self) -> &str {
        "auth"
    }

    fn router(&self) -> Router {
        let state = Arc::new(AuthState {
            db: self.db.clone(),
            jwt_secret: self.jwt_secret.clone(),
        });

        let user_provider = Arc::new(UserProvider::new(self.db.clone()));
        let role_provider = Arc::new(RoleProvider::new(self.db.clone()));

        Router::new()
            // auth endpoints
            .route("/auth/register", post(register))
            .route("/auth/login", post(login))
            .route("/auth/logout", post(logout))
            .route("/auth/me", axum::routing::get(me)
                .route_layer(middleware::from_fn(auth_middleware)))
            // CRUD routes (usable without admin module)
            .nest("/users", crud_router(user_provider))
            .nest("/roles", crud_router(role_provider))
            .with_state(state)
    }

    async fn migrate(&self, db: &DatabaseConnection) -> Result<(), AppError> {
        Migrator::up(db, None).await.map_err(|e| AppError::Internal(e.into()))
    }
}

// ── Request / Response types ──────────────────────────────────────────────────

#[derive(Deserialize)]
struct RegisterRequest {
    email: String,
    password: String,
    first_name: Option<String>,
    last_name: Option<String>,
}

#[derive(Deserialize)]
struct LoginRequest {
    email: String,
    password: String,
}

#[derive(Serialize)]
struct AuthResponse {
    token: String,
    user: user::Model,
}

// ── Handlers ──────────────────────────────────────────────────────────────────

async fn register(
    State(state): State<Arc<AuthState>>,
    jar: CookieJar,
    Json(req): Json<RegisterRequest>,
) -> Result<impl IntoResponse, AppError> {
    if req.email.is_empty() {
        return Err(AppError::BadRequest("Email is required".into()));
    }
    if req.password.len() < 8 {
        return Err(AppError::BadRequest("Password must be at least 8 characters".into()));
    }

    // check duplicate
    if user::Entity::find()
        .filter(user::Column::Email.eq(&req.email))
        .one(&state.db)
        .await?
        .is_some()
    {
        return Err(AppError::Conflict("Email already in use".into()));
    }

    // default "user" role
    let default_role = role::Entity::find()
        .filter(role::Column::Name.eq("user"))
        .one(&state.db)
        .await?
        .ok_or_else(|| AppError::Internal(anyhow::anyhow!("Default role not found")))?;

    let hash = bcrypt::hash(&req.password, bcrypt::DEFAULT_COST)
        .map_err(|e| AppError::Internal(e.into()))?;

    let now = Utc::now();
    let new_user = UserActiveModel {
        email: Set(req.email),
        password: Set(hash),
        first_name: Set(req.first_name),
        last_name: Set(req.last_name),
        role_id: Set(default_role.id),
        active: Set(true),
        created_at: Set(now),
        updated_at: Set(now),
        ..Default::default()
    };
    let user = new_user.insert(&state.db).await?;

    let token = generate_token(&user, &default_role.name, &state.jwt_secret)?;
    let jar = set_auth_cookie(jar, &token);

    Ok((StatusCode::CREATED, jar, Json(AuthResponse { token, user })))
}

async fn login(
    State(state): State<Arc<AuthState>>,
    jar: CookieJar,
    Json(req): Json<LoginRequest>,
) -> Result<impl IntoResponse, AppError> {
    if req.email.is_empty() || req.password.is_empty() {
        return Err(AppError::BadRequest("Email and password are required".into()));
    }

    let user = user::Entity::find()
        .filter(user::Column::Email.eq(&req.email))
        .one(&state.db)
        .await?
        .ok_or_else(|| AppError::Unauthorized("Invalid email or password".into()))?;

    if !user.active {
        return Err(AppError::Forbidden("Account is deactivated".into()));
    }

    let matches = bcrypt::verify(&req.password, &user.password)
        .map_err(|e| AppError::Internal(e.into()))?;
    if !matches {
        return Err(AppError::Unauthorized("Invalid email or password".into()));
    }

    let role = role::Entity::find_by_id(user.role_id)
        .one(&state.db)
        .await?
        .ok_or_else(|| AppError::Internal(anyhow::anyhow!("Role not found")))?;

    let token = generate_token(&user, &role.name, &state.jwt_secret)?;
    let jar = set_auth_cookie(jar, &token);

    Ok((jar, Json(AuthResponse { token, user })))
}

async fn logout(jar: CookieJar) -> impl IntoResponse {
    let jar = jar.remove(Cookie::build("auth_token").path("/"));
    (jar, Json(json!({ "message": "Logged out successfully" })))
}

async fn me(
    State(state): State<Arc<AuthState>>,
    axum::Extension(claims): axum::Extension<Claims>,
) -> Result<Json<user::Model>, AppError> {
    let user = user::Entity::find_by_id(claims.user_id)
        .one(&state.db)
        .await?
        .ok_or_else(|| AppError::NotFound("User not found".into()))?;
    Ok(Json(user))
}

// ── Middleware ────────────────────────────────────────────────────────────────

async fn auth_middleware(
    jar: CookieJar,
    mut req: axum::http::Request<axum::body::Body>,
    next: Next,
) -> Response {
    // try Authorization header first, then cookie
    let token = req
        .headers()
        .get(header::AUTHORIZATION)
        .and_then(|v| v.to_str().ok())
        .and_then(|v| v.strip_prefix("Bearer "))
        .map(str::to_owned)
        .or_else(|| jar.get("auth_token").map(|c| c.value().to_owned()));

    let secret = req
        .extensions()
        .get::<Arc<AuthState>>()
        .map(|s| s.jwt_secret.clone());

    // If we can't get secret from extensions, try state
    let claims = token.zip(secret).and_then(|(tok, sec)| {
        decode::<Claims>(&tok, &DecodingKey::from_secret(sec.as_bytes()), &Validation::default())
            .ok()
            .map(|d| d.claims)
    });

    match claims {
        Some(c) => {
            req.extensions_mut().insert(c);
            next.run(req).await
        }
        None => (
            StatusCode::UNAUTHORIZED,
            Json(json!({ "error": "Authentication required" })),
        )
            .into_response(),
    }
}

// ── JWT helpers ───────────────────────────────────────────────────────────────

fn generate_token(user: &user::Model, role: &str, secret: &str) -> Result<String, AppError> {
    let claims = Claims {
        user_id: user.id,
        email: user.email.clone(),
        role: role.to_owned(),
        exp: (Utc::now() + chrono::Duration::hours(24)).timestamp(),
        iat: Utc::now().timestamp(),
    };
    encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(secret.as_bytes()),
    )
    .map_err(|e| AppError::Internal(e.into()))
}

fn set_auth_cookie(jar: CookieJar, token: &str) -> CookieJar {
    let cookie = Cookie::build(("auth_token", token.to_owned()))
        .path("/")
        .http_only(true)
        .same_site(SameSite::Lax)
        .max_age(time::Duration::hours(24))
        .build();
    jar.add(cookie)
}
