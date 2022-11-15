-- +goose Up
-- +goose StatementBegin
create table records (
    id       bigserial,
    user_id  bigint,
    amount   numeric(19,4),
    category varchar(512),
    date     timestamp,
    
    primary key(id)
);

create index idx_records_user_id
on records(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table records;
-- +goose StatementEnd
