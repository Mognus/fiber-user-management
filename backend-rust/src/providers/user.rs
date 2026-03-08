use backend_core::{helpers, AppError, CrudResource, Field, FieldType, Filters, Schema};
use sea_orm::{DatabaseConnection, EntityTrait};

use crate::entities::user::{ActiveModel, Entity, Model};

pub struct UserProvider {
    db: DatabaseConnection,
}

impl UserProvider {
    pub fn new(db: DatabaseConnection) -> Self {
        Self { db }
    }
}

#[async_trait::async_trait]
impl CrudResource for UserProvider {
    type Entity = Entity;
    type ActiveModel = ActiveModel;

    fn db(&self) -> &DatabaseConnection {
        &self.db
    }

    fn schema(&self) -> Schema {
        Schema {
            name: "users".into(),
            display_name: "Users".into(),
            fields: vec![
                Field { name: "id".into(),         label: "ID".into(),         field_type: FieldType::Number,  readonly: true,  required: false, editable: true,  hidden: false, width: Some("80px".into()),  options: vec![] },
                Field { name: "email".into(),       label: "Email".into(),      field_type: FieldType::String,  readonly: false, required: true,  editable: true,  hidden: false, width: Some("200px".into()), options: vec![] },
                Field { name: "first_name".into(),  label: "First Name".into(), field_type: FieldType::String,  readonly: false, required: false, editable: true,  hidden: false, width: Some("140px".into()), options: vec![] },
                Field { name: "last_name".into(),   label: "Last Name".into(),  field_type: FieldType::String,  readonly: false, required: false, editable: true,  hidden: false, width: Some("140px".into()), options: vec![] },
                Field { name: "role_id".into(),     label: "Role".into(),       field_type: FieldType::Relation, readonly: false, required: true,  editable: true,  hidden: true,  width: None,                 options: vec![] },
                Field { name: "active".into(),      label: "Active".into(),     field_type: FieldType::Boolean, readonly: false, required: false, editable: true,  hidden: false, width: Some("100px".into()), options: vec![] },
                Field { name: "created_at".into(),  label: "Created".into(),    field_type: FieldType::Date,    readonly: true,  required: false, editable: false, hidden: false, width: Some("160px".into()), options: vec![] },
                Field { name: "updated_at".into(),  label: "Updated".into(),    field_type: FieldType::Date,    readonly: true,  required: false, editable: false, hidden: false, width: Some("160px".into()), options: vec![] },
            ],
            searchable: vec!["email".into(), "first_name".into(), "last_name".into()],
            filterable: vec!["role_id".into(), "active".into()],
        }
    }

    async fn get_by_id(&self, id: &str) -> Result<Model, AppError> {
        helpers::get_by_id::<Entity>(self.db(), id).await
    }

    async fn delete_by_id(&self, id: &str) -> Result<(), AppError> {
        helpers::delete_by_id::<Entity>(self.db(), id).await
    }
}
