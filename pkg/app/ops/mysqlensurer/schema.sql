--
-- Project table
--

CREATE TABLE IF NOT EXISTS Project (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED UNIQUE
) ENGINE=InnoDB;

--
-- Application table
--

CREATE TABLE IF NOT EXISTS Application (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  project_id VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  disabled BOOL GENERATED ALWAYS AS (data->>"$.disabled") STORED NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Command table
--

CREATE TABLE IF NOT EXISTS Command (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  project_id VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Deployment table
--

CREATE TABLE IF NOT EXISTS Deployment (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  project_id VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Environment table
--

CREATE TABLE IF NOT EXISTS Environment (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  project_id VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Piped table
--

CREATE TABLE IF NOT EXISTS Piped (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  project_id VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  disabled BOOL GENERATED ALWAYS AS (data->>"$.disabled") STORED NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- APIKey table
--

CREATE TABLE IF NOT EXISTS APIKey (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  project_id VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  disabled BOOL GENERATED ALWAYS AS (data->>"$.disabled") STORED NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;

--
-- Event table
--

CREATE TABLE IF NOT EXISTS Event (
  id BINARY(16) PRIMARY KEY,
  data JSON NOT NULL,
  project_id VARCHAR(50) GENERATED ALWAYS AS (data->>"$.project_id") STORED NOT NULL,
  created_at INT(11) GENERATED ALWAYS AS (data->>"$.created_at") STORED NOT NULL,
  updated_at INT(11) GENERATED ALWAYS AS (data->>"$.updated_at") STORED NOT NULL,
  extra VARCHAR(100) GENERATED ALWAYS AS (data->>"$._extra") STORED
) ENGINE=InnoDB;
