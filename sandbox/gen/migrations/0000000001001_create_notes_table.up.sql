create table notes
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    uuid       VARCHAR(36) UNIQUE,
    user_id    BIGINT,

    name       VARCHAR(255) NOT NULL,
    content    TEXT         NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id)
);