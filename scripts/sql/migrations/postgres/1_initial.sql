-- +migrate Up
CREATE TABLE "Organizations" (
	"OrganizationID" serial NOT NULL,
	"Name" VARCHAR(255) NOT NULL,
	"Picture" bytea,
	CONSTRAINT "Organizations_pk" PRIMARY KEY ("OrganizationID")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Users" (
	"UserID" serial NOT NULL,
	"Username" VARCHAR(255) NOT NULL,
	"Password" VARCHAR(255) NOT NULL,
	"Shownname" VARCHAR(255),
	"EMail" VARCHAR(255) NOT NULL,
	"SuperUser" BOOLEAN NOT NULL,
	CONSTRAINT "Users_pk" PRIMARY KEY ("UserID")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Members" (
	"SectionID" serial NOT NULL,
	"UserID" serial NOT NULL,
	"Rights" serial NOT NULL,
	CONSTRAINT "Members_pk" PRIMARY KEY ("SectionID", "UserID")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Events" (
	"EventID" serial NOT NULL,
	"OrganizationID" serial NOT NULL,
	"Name" VARCHAR(255) NOT NULL,
	"Address" VARCHAR(255),
	"Start" TIMESTAMP NOT NULL,
	"End" TIMESTAMP,
	"Creator" integer,
	CONSTRAINT "Events_pk" PRIMARY KEY ("EventID")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Attendees" (
	"EventID" integer NOT NULL,
	"UserID" integer NOT NULL,
	"Commitment" integer NOT NULL,
	"Comment" integer NOT NULL
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Comments" (
	"CommentID" serial NOT NULL,
	"EventID" integer NOT NULL,
	"UserID" serial NOT NULL,
	"Comment" VARCHAR(255) NOT NULL,
	CONSTRAINT "Comments_pk" PRIMARY KEY ("CommentID")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "Section" (
	"SectionID" serial NOT NULL,
	"OrganizationID" integer NOT NULL,
	"Name" VARCHAR(255) NOT NULL,
	CONSTRAINT "Section_pk" PRIMARY KEY ("SectionID")
) WITH (
  OIDS=FALSE
);





ALTER TABLE "Members" ADD CONSTRAINT "Members_fk0" FOREIGN KEY ("SectionID") REFERENCES "Section"("SectionID");
ALTER TABLE "Members" ADD CONSTRAINT "Members_fk1" FOREIGN KEY ("UserID") REFERENCES "Users"("UserID");

ALTER TABLE "Events" ADD CONSTRAINT "Events_fk0" FOREIGN KEY ("OrganizationID") REFERENCES "Organizations"("OrganizationID");
ALTER TABLE "Events" ADD CONSTRAINT "Events_fk1" FOREIGN KEY ("Creator") REFERENCES "Users"("UserID");

ALTER TABLE "Attendees" ADD CONSTRAINT "Attendees_fk0" FOREIGN KEY ("EventID") REFERENCES "Events"("EventID");
ALTER TABLE "Attendees" ADD CONSTRAINT "Attendees_fk1" FOREIGN KEY ("UserID") REFERENCES "Users"("UserID");

ALTER TABLE "Comments" ADD CONSTRAINT "Comments_fk0" FOREIGN KEY ("EventID") REFERENCES "Events"("EventID");
ALTER TABLE "Comments" ADD CONSTRAINT "Comments_fk1" FOREIGN KEY ("UserID") REFERENCES "Users"("UserID");

ALTER TABLE "Section" ADD CONSTRAINT "Section_fk0" FOREIGN KEY ("OrganizationID") REFERENCES "Organizations"("OrganizationID");

-- +migrate Down
DROP TABLE "Organizations";
DROP TABLE "Users";
DROP TABLE "Members";
DROP TABLE "Events";
DROP TABLE "Attendees";
DROP TABLE "Comments";
DROP TABLE "Section";
