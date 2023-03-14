--
-- Application table indexes
--

-- index on `Disabled` and `UpdatedAt` DESC
CREATE INDEX application_disabled_updated_at_desc ON Application (Disabled, UpdatedAt DESC);

-- index on `Name` ASC and `UpdatedAt` DESC
ALTER TABLE Application ADD COLUMN Name VARCHAR(50) GENERATED ALWAYS AS (data->>"$.name") VIRTUAL NOT NULL;
CREATE INDEX application_name_updated_at_desc ON Application (Name, UpdatedAt DESC);

-- index on `Deleted` and `CreatedAt` ASC
-- TODO: Reconsider make this Deleted column as STORED GENERATED COLUMN
ALTER TABLE Application ADD COLUMN Deleted BOOL GENERATED ALWAYS AS (IF(data->>"$.deleted" = 'true', True, False)) VIRTUAL NOT NULL;
CREATE INDEX application_deleted_created_at_asc ON Application (Deleted, CreatedAt);

-- index on `Kind` ASC and `UpdatedAt` DESC
ALTER TABLE Application ADD COLUMN Kind INT GENERATED ALWAYS AS (IFNULL(data->>"$.kind", 0)) VIRTUAL NOT NULL;
CREATE INDEX application_kind_updated_at_desc ON Application (Kind, UpdatedAt DESC);

-- index on `SyncState.Status` ASC and `UpdatedAt` DESC
ALTER TABLE Application ADD COLUMN SyncState_Status INT GENERATED ALWAYS AS (IFNULL(data->>"$.sync_state.status", 0)) VIRTUAL NOT NULL;
CREATE INDEX application_sync_state_updated_at_desc ON Application (SyncState_Status, UpdatedAt DESC);

-- index on `ProjectId` ASC and `UpdatedAt` DESC
CREATE INDEX application_project_id_updated_at_desc ON Application (ProjectId, UpdatedAt DESC);

-- index on `PipedId` ASC and `UpdatedAt` DESC
ALTER TABLE Application ADD COLUMN PipedId VARCHAR(36) GENERATED ALWAYS AS (data->>"$.piped_id") VIRTUAL NOT NULL;
CREATE INDEX application_piped_id_updated_at_desc ON Application (PipedId, UpdatedAt DESC);

-- TODO: Should remove that statement after few releases.
DROP INDEX application_piped_id ON Application;

--
-- Command table indexes
--

-- index on `Status` ASC and `CreatedAt` ASC
ALTER TABLE Command ADD COLUMN Status INT GENERATED ALWAYS AS (IFNULL(data->>"$.status", 0)) VIRTUAL NOT NULL;
CREATE INDEX command_status_created_at_asc ON Command (Status, CreatedAt);

-- index on `PipedId` ASC
ALTER TABLE Command ADD COLUMN PipedId VARCHAR(36) GENERATED ALWAYS AS (data->>"$.piped_id") VIRTUAL NOT NULL;
CREATE INDEX command_piped_id ON Command (PipedId);

--
-- Deployment table indexes
--

-- index on `ApplicationId` ASC and `UpdatedAt` DESC
ALTER TABLE Deployment ADD COLUMN ApplicationId VARCHAR(36) GENERATED ALWAYS AS (data->>"$.application_id") VIRTUAL NOT NULL;
CREATE INDEX deployment_application_id_updated_at_desc ON Deployment (ApplicationId, UpdatedAt DESC);

-- index on `ApplicationName` ASC and `UpdatedAt` DESC
ALTER TABLE Deployment ADD COLUMN ApplicationName VARCHAR(36) GENERATED ALWAYS AS (data->>"$.application_name") VIRTUAL NOT NULL;
CREATE INDEX deployment_application_name_updated_at_desc ON Deployment (ApplicationName, UpdatedAt DESC);

-- index on `ProjectId` ASC and `UpdatedAt` DESC
CREATE INDEX deployment_project_id_updated_at_desc ON Deployment (ProjectId, UpdatedAt DESC);

-- index on `Kind` ASC and `UpdatedAt` DESC
ALTER TABLE Deployment ADD COLUMN Kind INT GENERATED ALWAYS AS (IFNULL(data->>"$.kind", 0)) VIRTUAL NOT NULL;
CREATE INDEX deployment_kind_updated_at_desc ON Deployment (Kind, UpdatedAt DESC);

-- index on `Status` ASC and `UpdatedAt` DESC
ALTER TABLE Deployment ADD COLUMN Status INT GENERATED ALWAYS AS (IFNULL(data->>"$.status", 0)) VIRTUAL NOT NULL;
CREATE INDEX deployment_status_updated_at_desc ON Deployment (Status, UpdatedAt DESC);

-- index on `PipedId` ASC
ALTER TABLE Deployment ADD COLUMN PipedId VARCHAR(36) GENERATED ALWAYS AS (data->>"$.piped_id") VIRTUAL NOT NULL;
CREATE INDEX deployment_piped_id ON Deployment (PipedId);

-- index on `CompletedAt` DESC and `Id` ASC
ALTER TABLE Deployment ADD COLUMN CompletedAt INT(11) GENERATED ALWAYS AS (data->>"$.completed_at") VIRTUAL NULL;
CREATE INDEX deployment_completed_at_desc_id ON Deployment (CompletedAt DESC, Id);

-- index on `CompletedAt` ASC and `Id` ASC
CREATE INDEX deployment_completed_at_id_asc ON Deployment (CompletedAt, Id);

-- index on `DeploymentChainId` ASC and `UpdatedAt` DESC
ALTER TABLE Deployment ADD COLUMN DeploymentChainId VARCHAR(36) GENERATED ALWAYS AS (data->>"$.deployment_chain_id") VIRTUAL NOT NULL;
ALTER TABLE Deployment MODIFY DeploymentChainId VARCHAR(36) GENERATED ALWAYS AS (IFNULL(data->>"$.deployment_chain_id", "")) VIRTUAL NOT NULL;
CREATE INDEX deployment_chain_id_updated_at_desc ON Deployment (DeploymentChainId, UpdatedAt DESC);

--
-- Event table indexes
--

-- index on `ProjectId` ASC and `CreatedAt` ASC
CREATE INDEX event_project_id_created_at_asc ON Event (ProjectId, CreatedAt);
-- index on `ProjectId` ASC and `CreatedAt` DESC
CREATE INDEX event_project_id_created_at_desc ON Event (ProjectId, CreatedAt DESC);
-- index on `ProjectId` ASC and `UpdatedAt` DESC
CREATE INDEX event_project_id_updated_at_desc ON Event (ProjectId, UpdatedAt DESC);

-- index on `EventKey` ASC, `Name` ASC, `ProjectId` ASC and `CreatedAt` DESC
ALTER TABLE Event ADD COLUMN EventKey VARCHAR(64) GENERATED ALWAYS AS (data->>"$.event_key") VIRTUAL NOT NULL, ADD COLUMN Name VARCHAR(50) GENERATED ALWAYS AS (data->>"$.name") VIRTUAL NOT NULL;
CREATE INDEX event_event_key_name_project_id_created_at_desc ON Event (EventKey, Name, ProjectId, CreatedAt DESC);

-- index on `ProjectId` ASC, `Status` ASC, CreatedAt DESC
ALTER TABLE Event ADD COLUMN Status INT GENERATED ALWAYS AS (data->>"$.status") VIRTUAL NOT NULL;
ALTER TABLE Event MODIFY COLUMN Status INT GENERATED ALWAYS AS (IFNULL(data->>"$.status", 0)) VIRTUAL NOT NULL;
CREATE INDEX event_project_id_status_created_at_desc ON Event (ProjectId, Status, CreatedAt DESC);

-- index on `ProjectId` ASC, `Status` ASC, UpdatedAt DESC
ALTER TABLE Event ADD COLUMN Status INT GENERATED ALWAYS AS (data->>"$.status") VIRTUAL NOT NULL;
ALTER TABLE Event MODIFY COLUMN Status INT GENERATED ALWAYS AS (IFNULL(data->>"$.status", 0)) VIRTUAL NOT NULL;
CREATE INDEX event_project_id_status_updated_at_desc ON Event (ProjectId, Status, UpdatedAt DESC);

-- index on `Name` ASC, `ProjectId` ASC, `UpdatedAt` DESC
ALTER TABLE Event ADD COLUMN Name VARCHAR(50) GENERATED ALWAYS AS (data->>"$.name") VIRTUAL NOT NULL;
CREATE INDEX event_name_project_id_updated_at_desc ON Event (Name, ProjectId, UpdatedAt DESC);

-- index on `Name` ASC, `ProjectId` ASC, `Status` ASC, UpdatedAt DESC
ALTER TABLE Event ADD COLUMN Name VARCHAR(50) GENERATED ALWAYS AS (data->>"$.name") VIRTUAL NOT NULL, ADD COLUMN Status INT GENERATED ALWAYS AS (data->>"$.status") VIRTUAL NOT NULL;
ALTER TABLE Event MODIFY COLUMN Status INT GENERATED ALWAYS AS (IFNULL(data->>"$.status", 0)) VIRTUAL NOT NULL;
CREATE INDEX event_name_project_id_status_updated_at_desc ON Event (Name, ProjectId, Status, UpdatedAt DESC);

--
-- Piped table indexes
--

-- index on `ProjectId` ASC
CREATE INDEX piped_project_id_asc ON Piped (ProjectId);

--
-- DeploymentChain table indexes
--

-- index on `ProjectId` ASC and `UpdatedAt` DESC
CREATE INDEX deploymentchain_project_id_updated_at_desc ON DeploymentChain (ProjectId, UpdatedAt DESC);

-- index on `Status` ASC and `UpdatedAt` DESC
ALTER TABLE DeploymentChain ADD COLUMN Status INT GENERATED ALWAYS AS (IFNULL(data->>"$.status", 0)) VIRTUAL NOT NULL;
CREATE INDEX deploymentchain_status_updated_at_desc ON Deployment (Status, UpdatedAt DESC);
