--
-- Project table
--

CREATE TABLE IF NOT EXISTS Project (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- Application table
--

CREATE TABLE IF NOT EXISTS Application (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Disabled BOOL GENERATED ALWAYS AS (IF(data->>"$.disabled" = 'true', True, False)) STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- Command table
--

CREATE TABLE IF NOT EXISTS Command (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- Deployment table
--

CREATE TABLE IF NOT EXISTS Deployment (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- Piped table
--

CREATE TABLE IF NOT EXISTS Piped (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Disabled BOOL GENERATED ALWAYS AS (IF(data->>"$.disabled" = 'true', True, False)) STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- APIKey table
--

CREATE TABLE IF NOT EXISTS APIKey (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Disabled BOOL GENERATED ALWAYS AS (IF(data->>"$.disabled" = 'true', True, False)) STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- Event table
--

CREATE TABLE IF NOT EXISTS Event (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- DeploymentChain table
--

CREATE TABLE IF NOT EXISTS DeploymentChain (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;

--
-- DeploymentTrace table
--

CREATE TABLE IF NOT EXISTS DeploymentTrace (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL
) ENGINE=InnoDB;
