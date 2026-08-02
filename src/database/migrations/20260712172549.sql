-- Create "annotations" table
CREATE TABLE "public"."annotations" (
 "id" text NOT NULL,
 "text" text NULL,
 "bold" boolean NULL,
 "italic" boolean NULL,
 "strikethrough" boolean NULL,
 "underline" boolean NULL,
 "code" boolean NULL,
 "color" text NULL,
 "block_id" text NULL,
 PRIMARY KEY ("id")
);
-- Create "pages" table
CREATE TABLE "public"."pages" (
 "id" text NOT NULL,
 "title" text NULL,
 "page_type" text NOT NULL,
 "last_edited" timestamptz NOT NULL,
 "emoji_icon" text NULL,
 "icon_link" text NULL,
 PRIMARY KEY ("id")
);
-- Create "blocks" table
CREATE TABLE "public"."blocks" (
 "id" text NOT NULL,
 "full_text" text NULL,
 "type" text NULL,
 "payload" jsonb NULL,
 "hash" text NULL,
 "page_id" text NULL,
 PRIMARY KEY ("id"),
 CONSTRAINT "fk_pages_blocks" FOREIGN KEY ("page_id") REFERENCES "public"."pages" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "saves" table
CREATE TABLE "public"."saves" (
 "id" text NOT NULL,
 "last_save" timestamptz NULL,
 PRIMARY KEY ("id")
);
-- Create "save_pages" table
CREATE TABLE "public"."save_pages" (
 "save_id" text NOT NULL,
 "page_id" text NOT NULL,
 PRIMARY KEY ("save_id", "page_id"),
 CONSTRAINT "fk_save_pages_page" FOREIGN KEY ("page_id") REFERENCES "public"."pages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
 CONSTRAINT "fk_save_pages_save" FOREIGN KEY ("save_id") REFERENCES "public"."saves" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_save_page" to table: "save_pages"
CREATE UNIQUE INDEX "idx_save_page" ON "public"."save_pages" ("save_id", "page_id");
