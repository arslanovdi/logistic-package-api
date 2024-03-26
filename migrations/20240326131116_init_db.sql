-- +goose Up
-- +goose StatementBegin
create table if not exists package
(
    id      bigserial
        constraint package_pk
            primary key,
    weight  integer,
    title   varchar(50),
    created timestamp,
    updated timestamp,
    removed boolean
);

comment on table package is 'logistic package omp db';

comment on column package.weight is 'grams';

alter table package
    owner to logistic;

create table if not exists package_events
(
    id         bigserial
        constraint package_events_pk
            primary key,
    package_id bigint  not null,
    type       integer not null,
    status     integer not null,
    payload    jsonb,
    updated    timestamp
);

comment on column package_events.payload is 'model.Package object';
alter table package
    owner to logistic;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS package;
DROP TABLE IF EXISTS package_events;
-- +goose StatementEnd
