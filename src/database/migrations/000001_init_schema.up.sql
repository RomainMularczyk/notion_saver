CREATE TABLE "Page"(
    "id" VARCAHR NOT NULL,
    "title" VARCHAR NOT NULL,
    "page_type" VARCHAR NOT NULL,
    "last_edited" TIMESTAMP,
    "icon" VARCHAR,
    CONSTRAINT "id_pkey" PRIMARY KEY("id")
)

CREATE UNIQUE INDEX ON "title" ("index_title");
