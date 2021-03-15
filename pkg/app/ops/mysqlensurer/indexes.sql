-- TODO: Add all required indexes for PipeCD.
-- index on `updated_at` field on `Application` table in ASC direction
CREATE INDEX application_updated_at_asc ON Application (updated_at ASC);

-- index on `created_at` field on `Application` table in DESC direction
CREATE INDEX application_created_at_desc ON Application (created_at DESC);
