CREATE TABLE IF NOT EXISTS followers (
    user_id bigint NOT NULL,
    follower_id bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, follower_id),
    -- cascade on delete because this is not an important data to keep if the user is deleted
    -- this is hard delete
    -- mostly software use soft delete because if a client delete user then we can rollback the user(from the logs or somehow) and this metadeta(relation with other is also needed)
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE
)