-- +goose Up
-- +goose StatementBegin
create table package
(
    id        bigserial
        constraint package_pk
            primary key,
    weight    integer,
    title     varchar(50),
    createdat timestamp
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
