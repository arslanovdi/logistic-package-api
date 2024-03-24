-- +goose Up
-- +goose StatementBegin
create table IF NOT EXISTS package
(
    id        bigint
        constraint package_pk
            primary key,
    title     varchar(50) not null,
    weight    bigint,
    createdAt timestamp   not null
);

comment on table package is 'logistic package omp db';

comment on column package.weight is 'grams';

alter table package
    owner to logistic;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS package;
-- +goose StatementEnd
