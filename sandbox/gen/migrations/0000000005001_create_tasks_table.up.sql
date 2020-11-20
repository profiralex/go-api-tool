create table tasks
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    uuid       VARCHAR(36) UNIQUE,
    user_id    BIGINT,
    note_id    BIGINT,

    name       VARCHAR(255) NOT NULL,
    content    TEXT         NOT NULL,
    completed  BOOL         NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (note_id) REFERENCES notes (id)
);