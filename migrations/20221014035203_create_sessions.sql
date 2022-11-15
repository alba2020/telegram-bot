-- +goose Up
-- +goose StatementBegin
create table sessions (
    user_id bigint not null,
    currency varchar(32),
    month_limit numeric(19,4)
);

create index idx_sessions_user_id
on records(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sessions;
-- +goose StatementEnd
