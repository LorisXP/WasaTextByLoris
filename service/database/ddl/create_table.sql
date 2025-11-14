CREATE TABLE Users (
    userID INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE CHECK (length(name) BETWEEN 3 AND 15),
    photo TEXT CHECK(photo IS NULL OR length(photo) > 0)
);

CREATE TABLE Groups (
    groupID INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK (length(name) BETWEEN 3 AND 15),
    photo TEXT CHECK(photo IS NULL OR length(photo) > 0),
    adminID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE Members (
    groupID INTEGER NOT NULL REFERENCES Groups(groupID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    userID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY(groupID, userID)
);

CREATE TABLE Events (
    groupID INTEGER REFERENCES Groups(groupID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    dt_event TEXT NOT NULL DEFAULT(datetime('now')),
    type TEXT NOT NULL CHECK(type IN ('kick', 'leave', 'entered')),
    user INTEGER NOT NULL REFERENCES Users(userID),
    PRIMARY KEY(groupID, dt_event)
);

CREATE TABLE Contents (
    contentID INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK(type IN('text', 'gif', 'photo')),
    content TEXT NOT NULL CHECK(length(content) BETWEEN 1 AND 13981013)
);

CREATE TABLE Messages (
    messageID INTEGER PRIMARY KEY AUTOINCREMENT,
    conversationID INTEGER NOT NULL REFERENCES Conversations(conversationID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    senderID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    sent_at TEXT NOT NULL DEFAULT(datetime('now')),
    status TEXT NOT NULL CHECK(status IN('received', 'read')),
    type TEXT NOT NULL CHECK(type IN('standard', 'forward')),
    contentID INTEGER NOT NULL REFERENCES Contents(contentID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
CREATE INDEX idx_search_message ON Messages(conversationID,sent_at);

CREATE TABLE ConversationsUser (
    conversationID INTEGER PRIMARY KEY AUTOINCREMENT,
    user1ID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    user2ID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    lastMessageID INTEGER,
    CHECK (user1ID <> user2ID)
);
CREATE INDEX idx_conversation_user ON ConversationsUser(user1ID, user2ID);

CREATE TABLE MessagesUser (
    messageID INTEGER PRIMARY KEY AUTOINCREMENT,
    conversationID INTEGER NOT NULL REFERENCES ConversationsUser(conversationID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    senderID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    receiverID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    sent_at TEXT NOT NULL DEFAULT(datetime('now')),
    status TEXT NOT NULL CHECK(status IN('received', 'read')),
    type TEXT NOT NULL CHECK(type IN('standard', 'forward')),
    contentID INTEGER NOT NULL REFERENCES Contents(contentID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
CREATE INDEX idx_messages_user ON MessagesUser(conversationID, sent_at);


CREATE TABLE ConversationsGroup (
    conversationID INTEGER PRIMARY KEY AUTOINCREMENT,
    userID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    groupID INTEGER NOT NULL REFERENCES Groups(groupID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    lastMessageID INTEGER
);
CREATE INDEX idx_conversation_group ON ConversationsGroup(userID, groupID);

CREATE TABLE MessagesGroup (
    messageID INTEGER PRIMARY KEY AUTOINCREMENT,
    conversationID INTEGER NOT NULL REFERENCES ConversationsGroup(conversationID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    senderID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    sent_at TEXT NOT NULL DEFAULT(datetime('now')),
    status TEXT NOT NULL CHECK(status IN('received', 'read')),
    type TEXT NOT NULL CHECK(type IN('standard', 'forward')),
    contentID INTEGER NOT NULL REFERENCES Contents(contentID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
CREATE INDEX idx_messages_group ON MessagesGroup(conversationID, sent_at);


CREATE TABLE Comments (
    commentID INTEGER PRIMARY KEY AUTOINCREMENT,

    -- Può essere un messaggio user→user
    messageUserID INTEGER REFERENCES MessagesUser(messageID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    -- Oppure un messaggio user→group
    messageGroupID INTEGER REFERENCES MessagesGroup(messageID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    reaction TEXT NOT NULL CHECK(length(reaction) BETWEEN 1 AND 4),

    -- Esattamente uno dei due deve essere NOT NULL
    CHECK (
        (messageUserID IS NOT NULL AND messageGroupID IS NULL)
        OR
        (messageUserID IS NULL AND messageGroupID IS NOT NULL)
    )
);

CREATE INDEX idx_comments_user ON Comments(messageUserID);
CREATE INDEX idx_comments_group ON Comments(messageGroupID);



