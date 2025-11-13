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

CREATE TABLE Conversations (
    conversationID INTEGER PRIMARY KEY AUTOINCREMENT,
    OfUserID INTEGER NOT NULL REFERENCES Users(userID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    type TEXT NOT NULL CHECK(type IN('user', 'group')),
    receiverID INTEGER NOT NULL, /*userID or groupID*/
    lastMessageID INTEGER NOT NULL
);
CREATE INDEX idx_conversation_search ON Conversations(OfUserID, receiverID);


CREATE TABLE Comments (
    commentID INTEGER PRIMARY KEY AUTOINCREMENT,
    messageID INTEGER NOT NULL REFERENCES Messages(messageID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    reaction TEXT NOT NULL CHECK (length(reaction) BETWEEN 1 AND 4  -- evita stringhe lunghe)
);
CREATE INDEX idx_join_messages ON Comments(messageID);


