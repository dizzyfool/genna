drop schema if exists "public" cascade;
drop schema if exists "geo" cascade;

create schema "public";

create extension if not exists "uuid-ossp";

create table "projects"
(
    "projectId" uuid   not null default uuid_generate_v4(),
    "code"      uuid,
    "name"      text   not null,

    primary key ("projectId")
);

create table "users"
(
    "userId"    serial      not null,
    "email"     varchar(64) not null,
    "activated" bool        not null default false,
    "name"      varchar(128),
    "countryId" integer,
    "avatar"    bytea       not null,
    "avatarAlt" bytea,
    "apiKeys"   bytea[],
    "loggedAt"  timestamp,

    primary key ("userId")
);

create schema "geo";

create table geo."countries"
(
    "countryId" serial     not null,
    "code"      varchar(3) not null,
    "coords"    integer[],

    primary key ("countryId")
);

alter table "users"
    add constraint "fk_user_country"
        foreign key ("countryId")
            references geo."countries" ("countryId") on update restrict on delete restrict;