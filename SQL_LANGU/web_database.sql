SET TIME ZONE 'Asia/Shanghai';
DROP TRIGGER IF EXISTS alikes_insert_trigger ON ALikes CASCADE;
DROP TRIGGER IF EXISTS alikes_delete_trigger ON ALikes CASCADE;
DROP TRIGGER IF EXISTS clikes_insert_trigger ON CLikes CASCADE;
DROP TRIGGER IF EXISTS clikes_delete_trigger ON CLikes CASCADE;
DROP TRIGGER IF EXISTS acomments_insert_reply_trigger ON AComments CASCADE;
DROP TRIGGER IF EXISTS acomments_delete_reply_trigger ON AComments CASCADE;

DROP FUNCTION IF EXISTS update_article_like_count() CASCADE;
DROP FUNCTION IF EXISTS update_comment_like_count() CASCADE;
DROP FUNCTION IF EXISTS update_reply_count() CASCADE;

DROP FUNCTION IF EXISTS add_user_group(
    group_name VARCHAR(50)
) CASCADE;
DROP FUNCTION IF EXISTS add_user_min(
    p_gid INT,
    p_user_name VARCHAR(50),
    p_user_avatar CHAR(64)
);
DROP FUNCTION IF EXISTS add_article_min(
    p_author INT,
    p_title VARCHAR(128),
    p_description VARCHAR(512),
    p_body TEXT
);
DROP FUNCTION IF EXISTS comment_article(
    p_aid INT,
    p_uid INT,
    p_msg VARCHAR(500)
);
DROP FUNCTION IF EXISTS comment_reply(
    p_cmid INT,
    p_uid INT,
    p_msg VARCHAR(500)
);
DROP FUNCTION IF EXISTS comment_remove(
    p_cmid INT
);

/*
Std
*/
DROP FUNCTION IF EXISTS register_user(
    p_user_name VARCHAR(50),
    p_user_avatar CHAR(64),
    p_email VARCHAR(50),
    p_passwd CHAR(64)
);
DROP FUNCTION IF EXISTS check_authorization(
    p_auth CHAR(64)
);
DROP FUNCTION IF EXISTS file_upload(
    p_uid int,
    p_hash CHAR(64),
    p_filename VARCHAR(256),
    p_type VARCHAR(50)
);
DROP FUNCTION IF EXISTS check_max_online(
    p_uid INT
);
DROP FUNCTION IF EXISTS generate_random_string(
    length INT
);
DROP FUNCTION IF EXISTS get_timestamp_diff_s(
    p_time1 TIMESTAMP WITH TIME ZONE,
    p_time2 TIMESTAMP WITH TIME ZONE
);
DROP FUNCTION IF EXISTS reset_err(
    p_cid BIGINT,
    p_uid INT
);
DROP FUNCTION IF EXISTS login_account(
    p_user VARCHAR(50),
    p_passwd CHAR(64),
    p_ip INET,
    p_port INT
);
DROP FUNCTION IF EXISTS get_connection_id(
    p_ip INET,
    p_port INT
);
DROP FUNCTION IF EXISTS get_ip_connection_id(
    p_ip INET
);
DROP FUNCTION IF EXISTS get_ava_try(
    p_uid INT,
    p_ip INET
);
DROP FUNCTION IF EXISTS add_err(
    p_uid INT,
    p_ip INET
);
DROP FUNCTION IF EXISTS get_user_group(
    p_uid INT
);
DROP FUNCTION IF EXISTS get_control_mode(
    p_key VARCHAR(50)
);
DROP FUNCTION IF EXISTS get_uri_id(
    p_url VARCHAR(256)
);
DROP FUNCTION IF EXISTS get_url_ids(
    p_url VARCHAR(256)
);
DROP FUNCTION IF EXISTS url_acc_control(
    p_uid INT,
    p_url VARCHAR(256)
);
DROP FUNCTION IF EXISTS ip_acc_control(
    p_ip INET
);
DROP FUNCTION IF EXISTS acc_control(
    p_ip INET,
    p_url VARCHAR(256),
    p_uid INT
);
DROP FUNCTION IF EXISTS get_domain_id(
    p_domain VARCHAR(50)
);
DROP FUNCTION IF EXISTS add_inet_action(
    p_uid INT,
    p_ip INET,
    p_port INT,
    p_auth CHAR(64),
    p_domain VARCHAR(50),
    p_uri VARCHAR(256),
    p_http_code SMALLINT,
    p_server_code SMALLINT
);

DROP TABLE IF EXISTS Settings CASCADE;
DROP TABLE IF EXISTS GithubUsers CASCADE;
DROP TABLE IF EXISTS Users CASCADE;
DROP TABLE IF EXISTS LoginTypes CASCADE;
DROP TABLE IF EXISTS Authorizations CASCADE;
DROP TABLE IF EXISTS iNetActions CASCADE;
DROP TABLE IF EXISTS Connections CASCADE;
DROP TABLE IF EXISTS Domains CASCADE;
DROP TABLE IF EXISTS Articles CASCADE;
DROP TABLE IF EXISTS ADecorations CASCADE;
DROP TABLE IF EXISTS ALikes CASCADE;
DROP TABLE IF EXISTS Tags CASCADE;
DROP TABLE IF EXISTS ATags CASCADE;
DROP TABLE IF EXISTS AComments CASCADE;
DROP TABLE IF EXISTS CLikes CASCADE;
DROP TABLE IF EXISTS UserGroups CASCADE;
DROP TABLE IF EXISTS ArticleAccess CASCADE;
DROP TABLE IF EXISTS ArticleWhite CASCADE;
DROP TABLE IF EXISTS ArticleBlack CASCADE;
DROP TABLE IF EXISTS Urls CASCADE;
DROP TABLE IF EXISTS UrlWhite CASCADE;
DROP TABLE IF EXISTS UrlBlack CASCADE;
DROP TABLE IF EXISTS IPControl CASCADE;
DROP TABLE IF EXISTS LoginControl CASCADE;
DROP TABLE IF EXISTS ResourceTypes CASCADE;
DROP TABLE IF EXISTS Resources CASCADE;

/*
Settings Container
*/

CREATE TABLE IF NOT EXISTS Settings (
                                        sid SERIAL PRIMARY KEY,
                                        item_name VARCHAR(50) NOT NULL,
                                        item_val VARCHAR(50) NOT NULL
);

COMMENT ON COLUMN Settings.sid IS 'Setting item id';
COMMENT ON COLUMN Settings.item_name IS 'Setting item name';
COMMENT ON COLUMN Settings.item_val IS 'Setting item value';

/*
Users Container
*/

CREATE TABLE IF NOT EXISTS GithubUsers(
                                          github_id BIGINT PRIMARY KEY,
                                          name_val VARCHAR(50) NOT NULL,
                                          email VARCHAR(50)
);

COMMENT ON COLUMN GithubUsers.github_id IS 'Github account id';
COMMENT ON COLUMN GithubUsers.name_val IS 'Github user name';
COMMENT ON COLUMN GithubUsers.email IS 'Github account email address';

CREATE TABLE IF NOT EXISTS Users(
                                    uid SERIAL PRIMARY KEY,
                                    github_id BIGINT REFERENCES GithubUsers(github_id),
                                    gid INT NOT NULL,
                                    user_name VARCHAR(50) NOT NULL,
                                    user_avatar CHAR(64) NOT NULL,
                                    email VARCHAR(50),
                                    passwd CHAR(64),
                                    register_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN Users.uid IS 'User id';
COMMENT ON COLUMN Users.github_id IS 'Github account id';
COMMENT ON COLUMN Users.gid IS 'User group id';
COMMENT ON COLUMN Users.user_name IS 'User name';
COMMENT ON COLUMN Users.user_avatar IS 'User uri (hash)';
COMMENT ON COLUMN Users.email IS 'User email address';
COMMENT ON COLUMN Users.passwd IS 'User password (hash)';
COMMENT ON COLUMN Users.register_time IS 'Account register time';

CREATE TABLE IF NOT EXISTS LoginTypes(
                                         login_type_id SERIAL PRIMARY KEY,
                                         type_name VARCHAR(50) NOT NULL
);

COMMENT ON COLUMN LoginTypes.login_type_id IS 'Login type id';
COMMENT ON COLUMN LoginTypes.type_name IS 'Login type name';

CREATE TABLE IF NOT EXISTS Authorizations(
                                             auth CHAR(64) PRIMARY KEY,
                                             uid INT NOT NULL REFERENCES Users(uid),
                                             type_id INT NOT NULL REFERENCES LoginTypes(login_type_id),
                                             cid BIGINT NOT NULL,
                                             time_rc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN Authorizations.auth IS 'Authorization cookie value (hash)';
COMMENT ON COLUMN Authorizations.uid IS 'User id of the authorization cookie';
COMMENT ON COLUMN Authorizations.type_id IS 'Login type';
COMMENT ON COLUMN Authorizations.cid IS 'Login IP/Port';
COMMENT ON COLUMN Authorizations.time_rc IS 'Login time';

/*
Connections Container
*/

CREATE TABLE IF NOT EXISTS Connections(
                                          cid BIGSERIAL PRIMARY KEY,
                                          ip INET NOT NULL,
                                          port INT CHECK (port >= 0 AND port <= 65535)
);

COMMENT ON COLUMN Connections.cid IS 'Connection id';
COMMENT ON COLUMN Connections.ip IS 'Connection ip address';
COMMENT ON COLUMN Connections.port IS 'Connection port';

CREATE TABLE IF NOT EXISTS Domains(
                                      domain_id SERIAL PRIMARY KEY,
                                      domain_val VARCHAR(50) NOT NULL
);

COMMENT ON COLUMN Domains.domain_id IS 'Domain id';
COMMENT ON COLUMN Domains.domain_val IS 'Domain name/value';

CREATE TABLE IF NOT EXISTS iNetActions(
                                          naid SERIAL PRIMARY KEY,
                                          uid INT REFERENCES Users(uid),
                                          cid BIGINT NOT NULL REFERENCES Connections(cid),
                                          domain_id INT NOT NULL REFERENCES Domains(domain_id),
                                          uri_val VARCHAR(128) NOT NULL,
                                          auth CHAR(64),
                                          time_rc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                          ret_code SMALLINT,
                                          status_code SMALLINT
);

COMMENT ON COLUMN iNetActions.naid IS 'Network action id';
COMMENT ON COLUMN iNetActions.uid IS 'Request user id';
COMMENT ON COLUMN iNetActions.cid IS 'Connection id';
COMMENT ON COLUMN iNetActions.domain_id IS 'Target domain';
COMMENT ON COLUMN iNetActions.uri_val IS 'Request uri';
COMMENT ON COLUMN iNetActions.auth IS 'Authorization of user in this request';
COMMENT ON COLUMN iNetActions.time_rc IS 'Request TIMESTAMP WITH TIME ZONE';
COMMENT ON COLUMN iNetActions.ret_code IS 'Response code in http protoco';
COMMENT ON COLUMN iNetActions.status_code IS 'Response code which is customized';

/*
Artciles Container
*/

CREATE TABLE IF NOT EXISTS Articles(
                                       aid SERIAL PRIMARY KEY,
                                       author INT NOT NULL REFERENCES Users(uid),
                                       title VARCHAR(128) NOT NULL,
                                       description VARCHAR(512) NOT NULL,
                                       release_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       modify_time TIMESTAMP WITH TIME ZONE,
                                       like_count INT NOT NULL DEFAULT 0,
                                       body TEXT NOT NULL,
                                       cmt_count INT NOT NULL DEFAULT 0,
                                       view_count INT NOT NULL DEFAULT 0
);

COMMENT ON COLUMN Articles.aid IS 'Article id';
COMMENT ON COLUMN Articles.author IS 'Article author id';
COMMENT ON COLUMN Articles.title IS 'Article title';
COMMENT ON COLUMN Articles.description IS 'Article''s short description';
COMMENT ON COLUMN Articles.release_time IS 'Article release time';
COMMENT ON COLUMN Articles.modify_time IS 'Article modify time';
COMMENT ON COLUMN Articles.like_count IS 'Count of the click of ''like''';
COMMENT ON COLUMN Articles.body IS 'Article body';
COMMENT ON COLUMN Articles.cmt_count IS 'Count of comments under this article';
COMMENT ON COLUMN Articles.view_count IS 'Count of times that this article has been viewed';

CREATE TABLE IF NOT EXISTS ADecorations(
                                           aid INT PRIMARY KEY REFERENCES Articles(aid),
                                           cover CHAR(64),
                                           head CHAR(64),
                                           background CHAR(64)
);

COMMENT ON COLUMN ADecorations.aid IS 'Article id';
COMMENT ON COLUMN ADecorations.cover IS 'Article cover image  uri (hash)';
COMMENT ON COLUMN ADecorations.head IS 'Article head image uri (hash)';
COMMENT ON COLUMN ADecorations.background IS 'Article background image uri (hash)';

CREATE TABLE IF NOT EXISTS ALikes(
                                     uid INT REFERENCES Users(uid),
                                     aid INT REFERENCES Articles(aid),
                                     time_rc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (uid, aid)
);

COMMENT ON COLUMN ALikes.uid IS 'Who likes';
COMMENT ON COLUMN ALikes.aid IS 'Which article is liked';
COMMENT ON COLUMN ALikes.time_rc IS 'When likes';

CREATE TABLE IF NOT EXISTS Tags(
                                   tid SERIAL PRIMARY KEY,
                                   tag VARCHAR(50) NOT NULL
);

COMMENT ON COLUMN Tags.tid IS 'Tag id';
COMMENT ON COLUMN Tags.tag IS 'Tag name';

CREATE TABLE IF NOT EXISTS ATags(
                                    aid INT REFERENCES Articles(aid),
                                    tid INT REFERENCES Tags(tid),
                                    PRIMARY KEY (aid, tid)
);

COMMENT ON COLUMN ATags.aid IS 'Article id';
COMMENT ON COLUMN ATags.tid IS 'Tag id';

/*
Articles Container -> Comments Container
*/

CREATE TABLE IF NOT EXISTS AComments(
                                        cmid SERIAL PRIMARY KEY,
                                        sender INT NOT NULL REFERENCES Users(uid),
                                        aid INT REFERENCES Articles(aid),
                                        reply INT REFERENCES AComments(cmid),
                                        msg VARCHAR(500) NOT NULL,
                                        time_rc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        reply_count INT NOT NULL DEFAULT 0,
                                        like_count INT NOT NULL DEFAULT 0,
                                        removed BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON COLUMN AComments.cmid IS 'Comment id';
COMMENT ON COLUMN AComments.sender IS 'Comment sender id';
COMMENT ON COLUMN AComments.msg IS 'Comment message';
COMMENT ON COLUMN AComments.time_rc IS 'Comment send time';
COMMENT ON COLUMN AComments.reply_count IS 'Reply count of this comment';
COMMENT ON COLUMN AComments.like_count IS 'Like count of this comment';
COMMENT ON COLUMN AComments.reply IS 'Which comment is replied by this comment';
COMMENT ON COLUMN AComments.aid IS 'Which article is replied by this comment';
COMMENT ON COLUMN AComments.removed IS 'Is this comment removed?';

CREATE TABLE IF NOT EXISTS CLikes(
                                     uid INT REFERENCES Users(uid),
                                     cmid INT REFERENCES AComments(cmid),
                                     time_rc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (uid, cmid)
);

COMMENT ON COLUMN CLikes.uid IS 'Who likes';
COMMENT ON COLUMN CLikes.cmid IS 'Which comment is liked';
COMMENT ON COLUMN Clikes.time_rc IS 'When likes';

/*
Triggers
*/

CREATE OR REPLACE FUNCTION update_article_like_count()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE Articles
        SET like_count = like_count + 1
        WHERE aid = NEW.aid;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE Articles
        SET like_count = like_count - 1
        WHERE aid = OLD.aid;
        RETURN OLD;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_comment_like_count()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE AComments
        SET like_count = like_count + 1
        WHERE cmid = NEW.cmid;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE AComments
        SET like_count = like_count - 1
        WHERE cmid = OLD.cmid;
        RETURN OLD;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_reply_count()
    RETURNS TRIGGER AS
$BODY$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.reply IS NOT NULL THEN
            UPDATE AComments
            SET reply_count = reply_count + 1
            WHERE cmid = NEW.reply;
        ELSIF NEW.aid IS NOT NULL THEN
            UPDATE Articles
            SET cmt_count = cmt_count + 1
            WHERE aid = NEW.aid;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.reply IS NOT NULL THEN
            UPDATE AComments
            SET reply_count = reply_count - 1
            WHERE cmid = OLD.reply;
        ELSIF OLD.aid IS NOT NULL THEN
            UPDATE Articles
            SET cmt_count = cmt_count - 1
            WHERE aid = OLD.aid;
        END IF;
        RETURN OLD;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER alikes_insert_trigger
    AFTER INSERT ON ALikes
    FOR EACH ROW
EXECUTE FUNCTION update_article_like_count();

CREATE TRIGGER alikes_delete_trigger
    AFTER DELETE ON ALikes
    FOR EACH ROW
EXECUTE FUNCTION update_article_like_count();

CREATE TRIGGER clikes_insert_trigger
    AFTER INSERT ON CLikes
    FOR EACH ROW
EXECUTE FUNCTION update_comment_like_count();

CREATE TRIGGER clikes_delete_trigger
    AFTER DELETE ON CLikes
    FOR EACH ROW
EXECUTE FUNCTION update_comment_like_count();

CREATE TRIGGER acomments_insert_reply_trigger
    AFTER INSERT ON Acomments
    FOR EACH ROW
EXECUTE FUNCTION update_reply_count();

CREATE TRIGGER acomments_delete_reply_trigger
    AFTER DELETE ON Acomments
    FOR EACH ROW
EXECUTE FUNCTION update_reply_count();

/*
Security Container
*/

CREATE TABLE IF NOT EXISTS UserGroups(
                                         gid SERIAL PRIMARY KEY,
                                         name_val VARCHAR(50) NOT NULL
);

COMMENT ON COLUMN UserGroups.gid IS 'User group id';
COMMENT ON COLUMN UserGroups.name_val IS 'User group name';

CREATE TABLE IF NOT EXISTS ArticleAccess(
                                            aid INT PRIMARY KEY REFERENCES Articles(aid),
                                            whitelist BOOLEAN NOT NULL DEFAULT FALSE
);

COMMENT ON COLUMN ArticleAccess.aid IS 'Article id';
COMMENT ON COLUMN ArticleAccess.whitelist IS 'Is artcle under whitelist mode?';

CREATE TABLE IF NOT EXISTS ArticleWhite(
                                           aid INT REFERENCES Articles(aid),
                                           gid INT REFERENCES UserGroups(gid),
                                           PRIMARY KEY (aid, gid)
);

COMMENT ON COLUMN ArticleWhite.aid IS 'Whitelist article id';
COMMENT ON COLUMN ArticleWhite.gid IS 'Whitelist user group id';

CREATE TABLE IF NOT EXISTS ArticleBlack(
                                           aid INT REFERENCES Articles(aid),
                                           gid INT REFERENCES UserGroups(gid),
                                           PRIMARY KEY (aid, gid)
);

COMMENT ON COLUMN ArticleBlack.aid IS 'Blacklist article id';
COMMENT ON COLUMN ArticleBlack.gid IS 'Blacklist user group id';

CREATE TABLE IF NOT EXISTS Urls(
                                   url_id SERIAL PRIMARY KEY,
                                   whitelist BOOLEAN NOT NULL DEFAULT FALSE,
                                   url VARCHAR(256) NOT NULL
);

COMMENT ON COLUMN Urls.url_id IS 'Url id';
COMMENT ON COLUMN Urls.whitelist IS 'Whether Url control operates in whitelist mode';
COMMENT ON COLUMN Urls.url IS 'Url path value';

CREATE TABLE IF NOT EXISTS UrlWhite(
                                       url_id INT REFERENCES Urls(url_id),
                                       gid INT REFERENCES UserGroups(gid),
                                       PRIMARY KEY (url_id, gid)
);

COMMENT ON COLUMN UrlWhite.url_id IS 'Url id';
COMMENT ON COLUMN UrlWhite.gid IS 'Group id that permitted to access';

CREATE TABLE IF NOT EXISTS UrlBlack(
                                       url_id INT REFERENCES Urls(url_id),
                                       gid INT REFERENCES UserGroups(gid),
                                       PRIMARY KEY (url_id, gid)
);

COMMENT ON COLUMN UrlBlack.url_id IS 'Url id';
COMMENT ON COLUMN UrlBlack.gid IS 'Group id that is not allowed to access';

CREATE TABLE IF NOT EXISTS IPControl(
                                        ip INET,
                                        whitelist BOOLEAN DEFAULT FALSE,
                                        PRIMARY KEY (ip, whitelist)
);

COMMENT ON COLUMN IPControl.ip IS 'Target ip address';
COMMENT ON COLUMN IPControl.whitelist IS 'Control rule';

CREATE TABLE IF NOT EXISTS LoginControl(
                                           uid INT REFERENCES Users(uid),
                                           cid BIGINT REFERENCES Connections(cid),
                                           err_count SMALLINT NOT NULL DEFAULT 0,
                                           start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                           PRIMARY KEY (uid, cid)
);

COMMENT ON COLUMN LoginControl.uid IS 'User id who is tried to login';
COMMENT ON COLUMN LoginControl.cid IS 'Connection';
COMMENT ON COLUMN LoginControl.err_count IS 'Failure count';
COMMENT ON COLUMN LoginControl.start_time IS 'Time';

/*
Resources container
*/

CREATE TABLE IF NOT EXISTS ResourceTypes(
                                            resource_type_id SERIAL PRIMARY KEY,
                                            type_val VARCHAR(50) NOT NULL
);

COMMENT ON COLUMN ResourceTypes.resource_type_id IS 'Resource type id';
COMMENT ON COLUMN ResourceTypes.type_val IS 'Resource type name';

CREATE TABLE IF NOT EXISTS Resources(
                                        hash CHAR(64),
                                        name_val VARCHAR(256) NOT NULL,
                                        unique_val SMALLINT NOT NULL DEFAULT 0,
                                        type_id INT REFERENCES ResourceTypes(resource_type_id),
                                        uploader INT NOT NULL REFERENCES Users(uid),
                                        time_rc TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        removed BOOLEAN NOT NULL DEFAULT FALSE,
                                        PRIMARY KEY (hash, name_val)
);

COMMENT ON COLUMN Resources.hash IS 'Resource hash value';
COMMENT ON COLUMN Resources.name_val IS 'Resource filename';
COMMENT ON COLUMN Resources.unique_val IS 'Resource filename conflict solver';
COMMENT ON COLUMN Resources.type_id IS 'Resource''s type id';
COMMENT ON COLUMN Resources.uploader IS 'Resource''s uploader';
COMMENT ON COLUMN Resources.time_rc IS 'Resources''s upload time';
COMMENT ON COLUMN Resources.removed IS 'Whether resource is removed';

/*
Reference compete
*/

ALTER TABLE Users ADD FOREIGN KEY (gid) REFERENCES UserGroups(gid);
ALTER TABLE Authorizations ADD FOREIGN KEY (cid) REFERENCES Connections(cid);

/*
Functions
*/

CREATE OR REPLACE FUNCTION add_user_group(
    group_name VARCHAR(50)
)
    RETURNS TABLE (gid INT, name_val VARCHAR(50)) AS
$BODY$
DECLARE
    value_exists BOOLEAN;
BEGIN
    SELECT EXISTS(SELECT 1 FROM UserGroups WHERE UserGroups.name_val = group_name) INTO value_exists;
    IF NOT value_exists THEN
        INSERT INTO UserGroups(name_val) VALUES (group_name);
    END IF;
    RETURN QUERY SELECT * FROM UserGroups WHERE UserGroups.name_val = group_name;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_user_min(
    p_gid INT,
    p_user_name VARCHAR(50),
    p_user_avatar CHAR(64)
)
    RETURNS TABLE (uid INT, gid INT, user_name VARCHAR(50), user_avatar CHAR(64), register_time TIMESTAMP WITH TIME ZONE) AS
$BODY$
DECLARE
    value_exists BOOLEAN;
BEGIN
    SELECT EXISTS(SELECT 1 FROM Users WHERE Users.user_name = p_user_name) INTO value_exists;
    IF NOT value_exists THEN
        INSERT INTO Users(gid, user_name, user_avatar) VALUES (p_gid, p_user_name, p_user_avatar);
    END IF;
    RETURN QUERY SELECT Users.uid, Users.gid, Users.user_name, Users.user_avatar, Users.register_time FROM Users WHERE Users.user_name = p_user_name;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_article_min(
    p_author INT,
    p_title VARCHAR(128),
    p_description VARCHAR(512),
    p_body TEXT
)
    RETURNS TABLE (aid INT, title VARCHAR(128), description VARCHAR(512), body TEXT, release_time TIMESTAMP WITH TIME ZONE, like_count INT, cmt_count INT, view_count INT) AS
$BODY$
BEGIN
    INSERT INTO Articles(author, title, description, body)
    VALUES (p_author, p_title, p_description, p_body);
    RETURN QUERY SELECT Articles.aid, Articles.title, Articles.description, Articles.body, Articles.release_time, Articles.like_count, Articles.cmt_count, Articles.view_count
                 FROM Articles
                 WHERE Articles.author = p_author
                   AND Articles.title = p_title
                   AND Articles.description = p_description
                   AND Articles.body = p_body;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION comment_article(
    p_aid INT,
    p_uid INT,
    p_msg VARCHAR(500)
)
    RETURNS TABLE (cmid INT, sender INT, aid INT, msg VARCHAR(500),
                   time_rc TIMESTAMP WITH TIME ZONE, reply_count INT, like_count INT) AS
$BODY$
BEGIN
    INSERT INTO AComments(sender, aid, msg)
    VALUES (p_uid, p_aid, p_msg)
    RETURNING AComments.cmid, AComments.sender, AComments.aid, AComments.msg,
        AComments.time_rc, AComments.reply_count, AComments.like_count
        INTO cmid, sender, aid, msg,
            time_rc, reply_count, like_count;
    RETURN QUERY SELECT cmid, sender, aid, msg,
                        time_rc, reply_count, like_count;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION comment_reply(
    p_cmid INT,
    p_uid INT,
    p_msg VARCHAR(500)
)
    RETURNS TABLE (cmid INT, sender INT, reply INT, msg VARCHAR(500),
                   time_rc TIMESTAMP WITH TIME ZONE, reply_count INT, like_count INT) AS
$BODY$
BEGIN
    INSERT INTO AComments(sender, reply, msg)
    VALUES (p_uid, p_cmid, p_msg)
    RETURNING AComments.cmid, AComments.sender, AComments.reply, AComments.msg,
        AComments.time_rc, AComments.reply_count, AComments.like_count
        INTO cmid, sender, reply, msg,
            time_rc, reply_count, like_count;
    RETURN QUERY SELECT cmid, sender, reply, msg,
                        time_rc, reply_count, like_count;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION comment_remove(
    p_cmid INT
)
    RETURNS TABLE (cmid INT) AS
$BODY$
BEGIN
    UPDATE AComments
    SET removed = TRUE
    WHERE AComments.cmid = p_cmid
    RETURNING AComments.cmid INTO cmid;
    RETURN QUERY SELECT cmid;
END;
$BODY$ LANGUAGE plpgsql;

/*
Std
*/

CREATE OR REPLACE FUNCTION register_user(
    p_user_name VARCHAR(50),
    p_user_avatar CHAR(64),
    p_email VARCHAR(50),
    p_passwd CHAR(64)
)
    RETURNS TABLE (uid INT, msg TEXT) AS
$BODY$
DECLARE
    email_exist BOOLEAN;
    user_name_exist BOOLEAN;
BEGIN
    -- Check if p_email is used
    SELECT EXISTS(SELECT 1 FROM Users WHERE email = p_email) INTO email_exist;
    -- Check if p_user_name is used
    SELECT EXISTS(SELECT 1 FROM Users WHERE user_name = p_user_name) INTO user_name_exist;
    IF email_exist THEN
        RETURN QUERY SELECT -1, 'email_conflict';
    ELSIF user_name_exist THEN
        RETURN QUERY SELECT -1, 'user_name_conflict';
    ELSE
        INSERT INTO Users(gid, user_name, user_avatar, email, passwd) VALUES
            (2, p_user_name, p_user_avatar, p_email, p_passwd)
        RETURNING Users.uid INTO uid;
        RETURN QUERY SELECT uid, 'success';
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_authorization(
    p_auth CHAR(64)
)
    RETURNS TABLE (uid INT, group_name TEXT) AS
$BODY$
DECLARE
    p_uid INT;
    p_time TIMESTAMP WITH TIME ZONE;
    days_diff NUMERIC;
    group_name TEXT;
BEGIN
    BEGIN
        SELECT Authorizations.uid, Authorizations.time_rc INTO p_uid, p_time FROM Authorizations WHERE Authorizations.auth = p_auth;
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            p_uid := NULL;
            p_time := NULL;
    END;
    IF p_uid IS NULL THEN
        RETURN QUERY SELECT -1, 'anonymous';
    ELSE
        SELECT EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - p_time)) / (24 * 60 * 60) INTO days_diff;
        IF days_diff > 7 THEN
            DELETE FROM Authorizations WHERE Authorizations.auth = p_auth;
            RETURN QUERY SELECT -1, 'anonymous';
        ELSE
            SELECT UG.name_val INTO group_name FROM Users JOIN UserGroups UG on Users.gid = UG.gid WHERE Users.uid = p_uid;
            RETURN QUERY SELECT p_uid, group_name;
        END IF;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION file_upload(
    p_uid int,
    p_hash CHAR(64),
    p_filename VARCHAR(256),
    p_type VARCHAR(50)
)
    RETURNS TABLE (status INT, msg TEXT) AS
$BODY$
DECLARE
    t_type_id INT;
    t_hash_replace CHAR(64);
    t_unique_val INT;
    t_removed BOOLEAN;
BEGIN
    SELECT Resources.removed INTO t_removed FROM Resources WHERE Resources.hash = p_hash AND Resources.name_val = p_filename;
    IF t_removed IS NOT NULL THEN
        IF t_removed THEN
            UPDATE Resources SET removed = FALSE WHERE Resources.hash = p_hash AND Resources.name_val = p_filename;
            RETURN QUERY SELECT 0, 'success';
            RETURN;
        ELSE
            RETURN QUERY SELECT 0, 'exists';
            RETURN;
        END IF;
    END IF;
    IF p_type IS NOT NULL AND p_type != '' AND p_type != 'NULL' THEN
        SELECT ResourceTypes.resource_type_id INTO t_type_id FROM ResourceTypes WHERE ResourceTypes.type_val = p_type;
        IF t_type_id IS NULL THEN
            INSERT INTO ResourceTypes(type_val) VALUES (p_type) RETURNING ResourceTypes.resource_type_id INTO t_type_id;
        END IF;
    ELSE
        t_type_id := NULL;
    END IF;
    SELECT Resources.hash INTO t_hash_replace FROM Resources WHERE Resources.name_val = p_filename AND Resources.removed = TRUE ORDER BY Resources.unique_val LIMIT 1;
    IF t_hash_replace IS NOT NULL THEN
        UPDATE Resources SET hash = p_hash, removed = FALSE, type_id = t_type_id WHERE Resources.hash = t_hash_replace;
        IF NOT FOUND THEN
            RETURN QUERY SELECT -1, 'update failed';
            RETURN;
        END IF;
    ELSE
        SELECT MAX(Resources.unique_val) + 1 INTO t_unique_val FROM Resources WHERE Resources.name_val = p_filename;
        IF t_unique_val IS NOT NULL THEN
            INSERT INTO Resources(hash, name_val, type_id, uploader, unique_val) VALUES (p_hash, p_filename, t_type_id, p_uid, t_unique_val);
        ELSE
            INSERT INTO Resources(hash, name_val, type_id, uploader) VALUES (p_hash, p_filename, t_type_id, p_uid);
        END IF;
        IF NOT FOUND THEN
            RETURN QUERY SELECT -1, 'insert failed';
            RETURN;
        END IF;
    END IF;
    RETURN QUERY SELECT 0, 'success';
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_max_online(
    p_uid INT
)
    RETURNS VOID AS
$BODY$
DECLARE
    t_count INT;
    t_max_online INT;
BEGIN
    SELECT COUNT(Authorizations.auth) FROM Authorizations WHERE Authorizations.uid = p_uid INTO t_count;
    SELECT CAST(Settings.item_val AS INT) FROM Settings WHERE item_name = 'max_online' INTO t_max_online;
    IF t_count > t_max_online THEN
        DELETE FROM Authorizations
        WHERE Authorizations.auth IN (SELECT Au.auth FROM Authorizations AS Au WHERE Au.uid = p_uid ORDER BY Au.time_rc
                                      LIMIT (t_count - t_max_online));
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION generate_random_string(
    length INT
)
    RETURNS TEXT AS
$BODY$
DECLARE
    chars TEXT := '0123456789abcdef';
    random_string TEXT := '';
BEGIN
    WHILE length(random_string) < length LOOP
            random_string := random_string || substr(chars, floor(random() * length(chars))::INT, 1);
        END LOOP;
    RETURN random_string;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_timestamp_diff_s(
    p_time1 TIMESTAMP WITH TIME ZONE,
    p_time2 TIMESTAMP WITH TIME ZONE
)
    RETURNS INT AS
$BODY$
BEGIN
    RETURN EXTRACT(EPOCH FROM (p_time1 - p_time2));
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION reset_err(
    p_ip INET,
    p_uid INT
)
RETURNS VOID AS
$BODY$
DECLARE
    t_cid BIGINT;
BEGIN
    SELECT get_ip_connection_id(p_ip) INTO t_cid;
    IF EXISTS(SELECT 1 FROM LoginControl WHERE cid = t_cid AND uid = p_uid) THEN
        UPDATE LoginControl SET err_count = 0 WHERE cid = t_cid AND uid = p_uid;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION login_account(
    p_user VARCHAR(50),
    p_passwd CHAR(64),
    p_ip INET,
    p_port INT
)
    RETURNS TABLE (status INT, msg TEXT, auth TEXT) AS
$BODY$
DECLARE
    t_uid INT;
    t_passwd CHAR(64);
    t_au CHAR(64);
    t_login_type INT;
    t_insert_success BOOLEAN;
    t_cid BIGINT;
    t_ava BOOLEAN;
    t_left_err INT;
    t_left_s INT;
BEGIN
    SELECT Users.uid, Users.passwd INTO t_uid, t_passwd FROM Users WHERE CAST(Users.uid AS VARCHAR(50)) = p_user OR Users.user_name = p_user OR Users.email = p_user ORDER BY Users.uid
    LIMIT 1;
    IF t_uid IS NULL THEN
        RETURN QUERY SELECT -1, 'user not exists', '';
        RETURN;
    END IF;
    SELECT ava, left_err, left_s INTO t_ava, t_left_err, t_left_s FROM get_ava_try(t_uid, p_ip);
    IF NOT t_ava THEN
        RETURN QUERY SELECT -1, 'error limit exceeded', t_left_s::TEXT;
        RETURN;
    END IF;
    IF t_passwd IS NULL OR t_passwd != p_passwd THEN
        PERFORM add_err(t_uid, p_ip);
        RETURN QUERY SELECT -1, 'wrong password', (t_left_err - 1)::TEXT;
        RETURN;
    END IF;
    t_insert_success := FALSE;
    SELECT get_connection_id(p_ip, p_port) INTO t_cid;
    WHILE NOT t_insert_success LOOP
            BEGIN
                SELECT generate_random_string(64) INTO t_au;
                SELECT LoginTypes.login_type_id FROM LoginTypes WHERE type_name = 'password' INTO t_login_type;
                INSERT INTO Authorizations(auth, uid, type_id, cid) VALUES
                    (t_au, t_uid, t_login_type, t_cid);
                t_insert_success := TRUE;
            EXCEPTION
                WHEN UNIQUE_VIOLATION THEN
                    t_insert_success := FALSE;
            END;
        END LOOP;
    PERFORM check_max_online(t_uid);
    PERFORM reset_err(p_ip, t_uid);
    RETURN QUERY SELECT 0, 'success', CAST(t_au AS TEXT);
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_connection_id(
    p_ip INET,
    p_port INT
)
    RETURNS INT AS
$BODY$
DECLARE
    t_cid BIGINT;
BEGIN
    SELECT Connections.cid INTO t_cid FROM Connections WHERE Connections.ip = p_ip AND Connections.port = p_port;
    IF t_cid IS NULL THEN
        INSERT INTO Connections(ip, port) VALUES (p_ip, p_port) RETURNING cid INTO t_cid;
        RETURN t_cid;
    ELSE
        RETURN t_cid;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_ip_connection_id(
    p_ip INET
)
    RETURNS INT AS
$BODY$
DECLARE
    t_cid BIGINT;
BEGIN
    SELECT get_connection_id(p_ip, 0) INTO t_cid;
    RETURN t_cid;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_ava_try(
    p_uid INT,
    p_ip INET
)
    RETURNS TABLE (ava BOOLEAN, left_err INT, left_s INT) AS
$BODY$
DECLARE
    t_cid BIGINT;
    t_err_count SMALLINT;
    t_time TIMESTAMP WITH TIME ZONE;
    t_max_try SMALLINT;
    t_seconds INT;
    t_max_hours SMALLINT;
BEGIN
    SELECT get_ip_connection_id(p_ip) INTO t_cid;
    SELECT LoginControl.err_count, LoginControl.start_time INTO t_err_count, t_time FROM LoginControl WHERE LoginControl.uid = p_uid AND LoginControl.cid = t_cid;
    SELECT CAST(Settings.item_val AS SMALLINT) INTO t_max_try FROM Settings WHERE Settings.item_name = 'max_try';
    SELECT CAST(Settings.item_val AS SMALLINT) INTO t_max_hours FROM Settings WHERE Settings.item_name = 'try_hours';
    IF t_err_count IS NULL THEN
        INSERT INTO LoginControl(uid, cid) VALUES (p_uid, t_cid);
        RETURN QUERY SELECT TRUE, t_max_try::INT, 0;
        RETURN;
    ELSE
        SELECT get_timestamp_diff_s(CURRENT_TIMESTAMP, t_time) INTO t_seconds;
        IF t_seconds < t_max_hours * 3600 THEN
            RETURN QUERY SELECT (t_max_try - t_err_count) > 0, (t_max_try - t_err_count)::INT, t_max_hours * 3600 - t_seconds;
            RETURN;
        ELSE
            UPDATE LoginControl
            SET start_time = CURRENT_TIMESTAMP, err_count = 0
            WHERE LoginControl.uid = p_uid AND LoginControl.cid = t_cid;
            RETURN QUERY SELECT TRUE, t_max_try::INT, 0;
            RETURN;
        END IF;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_err(
    p_uid INT,
    p_ip INET
)
    RETURNS VOID AS
$BODY$
DECLARE
    t_cid BIGINT;
    t_err_count SMALLINT;
BEGIN
    SELECT get_ip_connection_id(p_ip) INTO t_cid;
    SELECT LoginControl.err_count INTO t_err_count FROM LoginControl WHERE LoginControl.uid = p_uid AND LoginControl.cid = t_cid;
    IF t_err_count IS NULL THEN
        INSERT INTO LoginControl(uid, cid, err_count) VALUES (p_uid, t_cid, 1);
        RETURN;
    ELSIF t_err_count = 0 THEN
        UPDATE LoginControl
        SET err_count = 1, start_time = CURRENT_TIMESTAMP
        WHERE LoginControl.uid = p_uid AND LoginControl.cid = t_cid;
        RETURN;
    ELSE
        UPDATE LoginControl
        SET err_count = t_err_count + 1
        WHERE LoginControl.uid = p_uid AND LoginControl.cid = t_cid;
        RETURN;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_user_group(
    p_uid INT
)
    RETURNS INT AS
$BODY$
DECLARE
    t_gid INT;
BEGIN
    SELECT Users.gid INTO t_gid FROM Users WHERE Users.uid = p_uid;
    IF t_gid IS NULL THEN
        SELECT UserGroups.gid INTO t_gid FROM UserGroups WHERE UserGroups.name_val = 'anonymous';
        RETURN t_gid;
    ELSE
        RETURN t_gid;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_control_mode(
    p_key VARCHAR(50)
)
    RETURNS INT AS
$BODY$
DECLARE
    t_val VARCHAR(50);
BEGIN
    SELECT Settings.item_val INTO t_val FROM Settings WHERE Settings.item_name = p_key;
    IF t_val IS NULL THEN
        RETURN 3;
    ELSIF t_val = 'blacklist' THEN
        RETURN 1;
    ELSIF t_val = 'whitelist' THEN
        RETURN 2;
    ELSE
        RETURN 3;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_uri_id(
    p_url VARCHAR(256)
)
    RETURNS INT AS
$BODY$
DECLARE
    t_url_id INT;
BEGIN
    SELECT Urls.url_id INTO t_url_id FROM Urls WHERE Urls.url = p_url;
    IF t_url_id IS NULL THEN
        INSERT INTO Urls(url) VALUES (p_url) RETURNING url_id INTO t_url_id;
        RETURN t_url_id;
    ELSE
        RETURN t_url_id;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_url_ids(
    p_url VARCHAR(256)
)
    RETURNS TABLE (id INT) AS
$BODY$
BEGIN
    RETURN QUERY SELECT Urls.url_id FROM Urls WHERE p_url ~ ('^' || Urls.url || '(\/|$)');
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION url_acc_control(
    p_uid INT,
    p_url VARCHAR(256)
)
    RETURNS BOOLEAN AS
$BODY$
DECLARE
    t_gid INT;
    t_mode INT;
    t_exists BOOLEAN;
    t_url_id INT;
BEGIN
    SELECT get_user_group(p_uid) INTO t_gid;
    SELECT get_control_mode('url_control') INTO t_mode;
    IF t_mode = 1 THEN
        t_exists := FALSE;
        FOR t_url_id IN (SELECT id FROM get_url_ids(p_url)) LOOP
            IF EXISTS(SELECT 1 FROM UrlBlack WHERE UrlBlack.url_id = t_url_id AND UrlBlack.gid = t_gid) THEN
                t_exists := TRUE;
                EXIT;
            END IF;
        END LOOP;
        IF t_exists THEN
            RETURN FALSE;
        ELSE
            RETURN TRUE;
        END IF;
    ELSIF t_mode = 2 THEN
        t_exists := FALSE;
        FOR t_url_id IN (SELECT id FROM get_url_ids(p_url)) LOOP
                IF EXISTS(SELECT 1 FROM UrlWhite WHERE UrlWhite.url_id = t_url_id AND UrlWhite.gid = t_gid) THEN
                    t_exists := TRUE;
                    EXIT;
                END IF;
            END LOOP;
        IF t_exists THEN
            RETURN TRUE;
        ELSE
            RETURN FALSE;
        END IF;
    ELSE
        RETURN TRUE;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ip_acc_control(
    p_ip INET
)
    RETURNS BOOLEAN AS
$BODY$
DECLARE
    t_mode INT;
BEGIN
    SELECT get_control_mode('ip_control') INTO t_mode;
    IF t_mode = 1 THEN
        IF EXISTS(SELECT 1 FROM IPControl WHERE IPControl.ip = p_ip AND IPControl.whitelist = FALSE) THEN
            RETURN FALSE;
        ELSE
            RETURN TRUE;
        END IF;
    ELSIF t_mode = 2 THEN
        IF EXISTS(SELECT 1 FROM IPControl WHERE IPControl.ip = p_ip AND IPControl.whitelist = TRUE) THEN
            RETURN TRUE;
        ELSE
            RETURN FALSE;
        END IF;
    ELSE
        RETURN TRUE;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION acc_control(
    p_ip INET,
    p_url VARCHAR(256),
    p_uid INT
)
    RETURNS BOOLEAN AS
$BODY$
BEGIN
    IF NOT ip_acc_control(p_ip) THEN
        RETURN FALSE;
    ELSIF NOT url_acc_control(p_uid, p_url) THEN
        RETURN FALSE;
    ELSE
        RETURN TRUE;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_domain_id(
    p_domain VARCHAR(50)
)
    RETURNS INT AS
$BODY$
DECLARE
    t_domain_id INT;
BEGIN
    SELECT Domains.domain_id INTO t_domain_id FROM Domains WHERE Domains.domain_val = p_domain;
    IF t_domain_id IS NULL THEN
        INSERT INTO Domains(domain_val) VALUES (p_domain) RETURNING Domains.domain_id INTO t_domain_id;
        RETURN t_domain_id;
    ELSE
        RETURN t_domain_id;
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_inet_action(
    p_uid INT,
    p_ip INET,
    p_port INT,
    p_auth CHAR(64),
    p_domain VARCHAR(50),
    p_uri VARCHAR(256),
    p_http_code SMALLINT,
    p_server_code SMALLINT
)
    RETURNS VOID AS
$BODY$
DECLARE
    t_cid BIGINT;
    t_domain_id INT;
BEGIN
    SELECT get_connection_id(p_ip, p_port) INTO t_cid;
    SELECT get_domain_id(p_domain) INTO t_domain_id;
    INSERT INTO iNetActions(uid, cid, domain_id, auth, uri_val, ret_code, status_code) VALUES (p_uid, t_cid, t_domain_id, p_auth, p_uri, p_http_code, p_server_code);
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_update_github_info(
    p_github_id BIGINT,
    p_github_username VARCHAR(50),
    p_github_email VARCHAR(50)
)
    RETURNS VOID AS
$BODY$
BEGIN
    IF EXISTS(SELECT 1 FROM GithubUsers WHERE GithubUsers.github_id = p_github_id) THEN
        UPDATE GithubUsers
            SET name_val = p_github_username, email = p_github_email
                WHERE github_id = p_github_id;
    ELSE
        INSERT INTO GithubUsers(github_id, name_val, email) VALUES
            (p_github_id, p_github_username, p_github_email);
    END IF;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION name_push_random(
    p_name TEXT
)
RETURNS TEXT AS
$BODY$
DECLARE
    c_random_chars TEXT := '0123456789abcdef';
BEGIN
    RETURN p_name || SUBSTR(c_random_chars, FLOOR(RANDOM() * LENGTH(c_random_chars)) , 1);
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION register_login_from_github(
    p_github_id BIGINT,
    p_github_username VARCHAR(50),
    p_github_email VARCHAR(50),
    p_avatar CHAR(64)
)
    RETURNS TABLE (status integer, msg text, auth text, uid INT) AS
$BODY$
DECLARE
    t_uid INT;
    t_name_raw VARCHAR(50) := p_github_username;
    t_name_conflict VARCHAR(50);
BEGIN
    PERFORM add_update_github_info(p_github_id, p_github_username, p_github_email);
    SELECT Users.uid INTO t_uid FROM Users WHERE github_id = p_github_id;
    IF t_uid IS NULL THEN
        -- user not exists
        -- detect user_name conflict
        IF EXISTS(SELECT 1 FROM Users WHERE USers.user_name = p_github_username) THEN
            IF LENGTH(t_name_raw) > 40 THEN
                t_name_raw = SUBSTR(t_name_raw, 0, 40);
            END IF;
            t_name_conflict := t_name_raw;
            -- name conflict
            WHILE TRUE LOOP
            BEGIN
                -- name generate
                IF LENGTH(t_name_conflict) >= 50 THEN
                    t_name_conflict := t_name_raw;
                END IF;
                SELECT name_push_random(t_name_conflict::TEXT) INTO t_name_conflict;
                IF NOT EXISTS(SELECT 1 FROM Users WHERE Users.user_name = t_name_conflict) THEN
                    EXIT;
                END IF;
            END;
            END LOOP;
            -- email conflict
            IF EXISTS(SELECT 1 FROM Users WHERE email = p_github_email) THEN
                INSERT INTO Users(github_id, gid, user_name, user_avatar)
                VALUES (p_github_id, 2, t_name_conflict, p_avatar) RETURNING Users.uid
                    INTO t_uid;
            ELSE
                INSERT INTO Users(github_id, gid, user_name, user_avatar, email)
                    VALUES (p_github_id, 2, t_name_conflict, p_avatar, p_github_email) RETURNING Users.uid
                        INTO t_uid;
            END IF;
        END IF;
    END IF;
    -- get uid
    -- to do
END;
$BODY$ LANGUAGE plpgsql;

/*
Run once
*/
INSERT INTO LoginTypes(type_name) VALUES ('password'), ('github');

INSERT INTO Settings(item_name, item_val) VALUES
                                              ('ip_control', 'blacklist'),
                                              ('url_control', 'blacklist'),
                                              ('article_control', 'blacklist'),
                                              ('max_days', '7'),
                                              ('max_online', '5'),
                                              ('max_try', '5'),
                                              ('try_hours', '5');
-- 1
SELECT * FROM add_user_group('admin');
-- 2
SELECT * FROM add_user_group('user');
-- 3
SELECT * FROM add_user_group('anonymous');
-- 1
SELECT * FROM add_user_min(1, 'Admin', '');
-- 2
SELECT * FROM add_user_min(3, 'Anonymous', '');
