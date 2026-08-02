-- Modify "saves" table
ALTER TABLE "public"."saves" ADD COLUMN "status" text NULL;
-- Drop "annotations" table
DROP TABLE "public"."annotations";
