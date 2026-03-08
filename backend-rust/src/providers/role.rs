use backend_core::{helpers, AppError, CrudResource, Field, FieldType, Filters, Schema};
use sea_orm::{DatabaseConnection, EntityTrait};

use crate::entities::role::{ActiveModel, Entity, Model};

pub struct RoleProvider {
    db: DatabaseConnection,
}

impl RoleProvider {
    pub fn new(db: DatabaseConnection) -> Self {
        Self { db }
    }
}

#[async_trait::async_trait]
impl CrudResource for RoleProvider {
    type Entity = Entity;
    type ActiveModel = ActiveModel;

    fn db(&self) -> &DatabaseConnection {
        &self.db
    }

    fn schema(&self) -> Schema {
        Schema {
            name: "roles".into(),
            display_name: "Roles".into(),
            fields: vec![
                Field { name: "id".into(),         label: "ID".into(),      field_type: FieldType::Number,  readonly: true,  required: false, editable: true,  hidden: false, width: Some("80px".into()),  options: vec![] },
                Field { name: "name".into(),        label: "Name".into(),    field_type: FieldType::String,  readonly: false, required: true,  editable: true,  hidden: false, width: Some("250px".into()), options: vec![] },
                Field { name: "created_at".into(),  label: "Created".into(), field_type: FieldType::Date,    readonly: true,  required: false, editable: false, hidden: false, width: Some("200px".into()), options: vec![] },
                Field { name: "updated_at".into(),  label: "Updated".into(), field_type: FieldType::Date,    readonly: true,  required: false, editable: false, hidden: false, width: Some("200px".into()), options: vec![] },
            ],
            searchable: vec!["name".into()],
            filterable: vec![],
        }
    }

    async fn get_by_id(&self, id: &str) -> Result<Model, AppError> {
        helpers::get_by_id::<Entity>(self.db(), id).await
    }

    async fn delete_by_id(&self, id: &str) -> Result<(), AppError> {
        helpers::delete_by_id::<Entity>(self.db(), id).await
    }
}
