-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create type activity_type_enum as enum (
    'Walking',
    'Yoga',
    'Stretching',
    'Cycling',
    'Swimming',
    'Dancing',
    'Hiking',
    'Running',
    'HIIT',
    'JumpRope'
);

create table if not exists activities (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    activity_type activity_type_enum,
    done_at timestamptz default null,
    duration_in_minutes int default 1,
    calories_burned int,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp,
    is_active bit default B'1'
);

CREATE INDEX IF NOT EXISTS activities_user_id_idx ON activities USING HASH (user_id);
CREATE INDEX IF NOT EXISTS activities_done_at_idx ON activities(user_id, done_at);
CREATE INDEX IF NOT EXISTS activities_calories_burned_idx ON activities(user_id, calories_burned);
CREATE INDEX IF NOT EXISTS activities_activity_type_idx ON activities(user_id, activity_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists activities;
drop type if exists activity_type_enum;
-- +goose StatementEnd
