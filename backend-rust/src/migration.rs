use sea_orm_migration::prelude::*;

pub struct Migrator;

#[async_trait::async_trait]
impl MigratorTrait for Migrator {
    fn migrations() -> Vec<Box<dyn MigrationTrait>> {
        vec![Box::new(m20240101_000001_create_roles_users::Migration)]
    }
}

mod m20240101_000001_create_roles_users {
    use sea_orm_migration::prelude::*;

    pub struct Migration;

    impl MigrationName for Migration {
        fn name(&self) -> &str {
            "m20240101_000001_create_roles_users"
        }
    }

    #[async_trait::async_trait]
    impl MigrationTrait for Migration {
        async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
            // roles table
            manager
                .create_table(
                    Table::create()
                        .table(Roles::Table)
                        .if_not_exists()
                        .col(ColumnDef::new(Roles::Id).big_integer().not_null().auto_increment().primary_key())
                        .col(ColumnDef::new(Roles::Name).string_len(50).not_null().unique_key())
                        .col(ColumnDef::new(Roles::CreatedAt).timestamp_with_time_zone().not_null())
                        .col(ColumnDef::new(Roles::UpdatedAt).timestamp_with_time_zone().not_null())
                        .to_owned(),
                )
                .await?;

            // users table
            manager
                .create_table(
                    Table::create()
                        .table(Users::Table)
                        .if_not_exists()
                        .col(ColumnDef::new(Users::Id).big_integer().not_null().auto_increment().primary_key())
                        .col(ColumnDef::new(Users::Email).string_len(255).not_null().unique_key())
                        .col(ColumnDef::new(Users::Password).string_len(255).not_null())
                        .col(ColumnDef::new(Users::FirstName).string_len(100))
                        .col(ColumnDef::new(Users::LastName).string_len(100))
                        .col(ColumnDef::new(Users::RoleId).big_integer().not_null())
                        .col(ColumnDef::new(Users::Active).boolean().not_null().default(true))
                        .col(ColumnDef::new(Users::CreatedAt).timestamp_with_time_zone().not_null())
                        .col(ColumnDef::new(Users::UpdatedAt).timestamp_with_time_zone().not_null())
                        .foreign_key(
                            ForeignKey::create()
                                .from(Users::Table, Users::RoleId)
                                .to(Roles::Table, Roles::Id)
                                .on_update(ForeignKeyAction::Cascade)
                                .on_delete(ForeignKeyAction::Restrict),
                        )
                        .to_owned(),
                )
                .await?;

            // seed default roles
            let now = "NOW()";
            manager
                .exec_stmt(
                    Query::insert()
                        .into_table(Roles::Table)
                        .columns([Roles::Name, Roles::CreatedAt, Roles::UpdatedAt])
                        .values_panic(["admin".into(), Expr::cust(now), Expr::cust(now)])
                        .values_panic(["user".into(), Expr::cust(now), Expr::cust(now)])
                        .values_panic(["guest".into(), Expr::cust(now), Expr::cust(now)])
                        .on_conflict(OnConflict::column(Roles::Name).do_nothing().to_owned())
                        .to_owned(),
                )
                .await
        }

        async fn down(&self, manager: &SchemaManager) -> Result<(), DbErr> {
            manager.drop_table(Table::drop().table(Users::Table).to_owned()).await?;
            manager.drop_table(Table::drop().table(Roles::Table).to_owned()).await
        }
    }

    #[derive(Iden)]
    enum Roles {
        Table,
        Id,
        Name,
        CreatedAt,
        UpdatedAt,
    }

    #[derive(Iden)]
    enum Users {
        Table,
        Id,
        Email,
        Password,
        FirstName,
        LastName,
        RoleId,
        Active,
        CreatedAt,
        UpdatedAt,
    }
}
