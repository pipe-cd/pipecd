--
-- Project table
--

CREATE TABLE IF NOT EXISTS Project (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED UNIQUE
) ENGINE=InnoDB;

--
-- Application table
--

CREATE TABLE IF NOT EXISTS Application (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Disabled BOOL GENERATED ALWAYS AS (data->>"$.disabled") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Command table
--

CREATE TABLE IF NOT EXISTS Command (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Deployment table
--

CREATE TABLE IF NOT EXISTS Deployment (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Environment table
--

CREATE TABLE IF NOT EXISTS Environment (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Piped table
--

CREATE TABLE IF NOT EXISTS Piped (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Disabled BOOL GENERATED ALWAYS AS (data->>"$.disabled") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- APIKey table
--

CREATE TABLE IF NOT EXISTS APIKey (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  Disabled BOOL GENERATED ALWAYS AS (data->>"$.disabled") STORED,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Event table
--

CREATE TABLE IF NOT EXISTS Event (
  Id BINARY(16) PRIMARY KEY,
  Data JSON NOT NULL,
  ProjectId VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  CreatedAt INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  UpdatedAt INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  Extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;
