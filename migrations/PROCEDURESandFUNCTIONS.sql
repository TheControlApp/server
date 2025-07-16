-- Some T-SQL functions do not have direct equivalents in PostgreSQL.
-- The following custom function is created to replicate the behavior of T-SQL's ISNUMERIC().
CREATE OR REPLACE FUNCTION isnumeric(text) RETURNS BOOLEAN AS $$
BEGIN
    PERFORM $1::numeric;
    RETURN TRUE;
EXCEPTION WHEN OTHERS THEN
    RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

-- Converted Stored Procedures and Views

CREATE OR REPLACE PROCEDURE USP_DeleteOutstanding(IN userID int DEFAULT 0)
AS $$
BEGIN
    DELETE FROM ControlAppCmd
    WHERE SubId = userID;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetChatLog(
    request_id int,
    other_id int,
    other_user varchar(255) DEFAULT ''
)
RETURNS TABLE (
    "Type" char(1),
    ChatTxt varchar,
    Date_Add timestamp
)
AS $$
DECLARE
    newother int;
BEGIN
    IF other_user != '' THEN
        SELECT "ID" INTO newother
        FROM Users
        WHERE "Screen Name" = other_user;
    ELSE
        newother := other_id;
    END IF;

    RETURN QUERY
    SELECT * FROM (
        SELECT 'S' AS "Type", c.ChatTxt, c.Date_Add
        FROM ChatLog c
        WHERE c.SenderID = request_id
        AND c.ReceiverID = newother

        UNION

        SELECT 'R' AS "Type", c.ChatTxt, c.Date_Add
        FROM ChatLog c
        WHERE c.ReceiverID = request_id
        AND c.SenderID = newother
    ) AS sub
    ORDER BY sub.Date_Add;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetDom(p_SubID int)
RETURNS TABLE (
    UserName varchar,
    DomID int
)
AS $$
BEGIN
    RETURN QUERY
    SELECT u."Screen Name" AS UserName, r.DomID
    FROM Relationship r
    JOIN Users u ON r.DomID = u.Id
    WHERE r.SubID = p_SubID;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetAppContent(IN userID int DEFAULT 0)
RETURNS TABLE (
    SenderId varchar(100),
    SenderName varchar,
    Content text
)
AS $$
BEGIN
    DELETE FROM ControlAppCmd
    WHERE SenderId IN (SELECT BlockeeId FROM Block WHERE BlockerId = userID)
      AND SubId = userID;

    DELETE FROM ControlAppCmd
    WHERE SenderId = -1
      AND SubId = userID;

    RETURN QUERY
    SELECT
        CAST(cac.SenderId AS varchar(100)),
        u."Screen Name" AS SenderName,
        cl.Content
    FROM ControlAppCmd cac
    JOIN CommandList cl ON cac.CmdId = cl.CmdId
    JOIN Users u ON cac.SenderId = u.Id
    WHERE cac.SubId = userID
    ORDER BY SendDate
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetAllAppContent(IN userID int DEFAULT 0)
RETURNS TABLE (
    CacId int,
    SenderId varchar(100),
    Content text,
    IsDom char(1)
)
AS $$
DECLARE
    Anon bit;
BEGIN
    DELETE FROM ControlAppCmd
    WHERE SenderId IN (SELECT BlockeeId FROM Block WHERE BlockerId = userID)
      AND SubId = userID;

    SELECT AnonCmd INTO Anon FROM Users WHERE ID = userID;

    IF Anon = B'0' THEN
        DELETE FROM ControlAppCmd
        WHERE SenderId = -1
          AND SubId = userID;
    END IF;

    RETURN QUERY
    SELECT
        cac.Id AS CacId,
        CAST(cac.SenderId AS varchar(100)) AS SenderId,
        cl.Content,
        CASE WHEN r.Id IS NULL THEN 'N' ELSE 'D' END AS IsDom
    FROM ControlAppCmd cac
    JOIN CommandList cl ON cac.CmdId = cl.CmdId
    LEFT JOIN Relationship r ON r.DomID = cac.SenderId AND r.SubID = userID
    WHERE cac.SubId = userID
    ORDER BY SendDate;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_DeleteInvite(
    IN p_SubId int,
    IN p_DomName varchar(250)
)
AS $$
BEGIN
    DELETE FROM Invites i
    USING Users u
    WHERE i.DomId = u.Id
      AND u."Screen Name" = p_DomName
      AND i.SubId = p_SubId;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_CmdComplete(IN userID int DEFAULT 0)
AS $$
DECLARE
    cmd_to_delete int;
BEGIN
    SELECT id INTO cmd_to_delete
    FROM ControlAppCmd c
    JOIN CommandList cl ON c.CmdId = cl.CmdId
    WHERE SubId = userID
    ORDER BY SendDate
    LIMIT 1;

    DELETE FROM ControlAppCmd
    WHERE id = cmd_to_delete;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION Usp_checkuser(IN p_ID int)
RETURNS int AS $$
DECLARE
    verified_status int;
BEGIN
    SELECT CASE WHEN vr.SubID IS NULL THEN u.Varified ELSE -1 END
    INTO verified_status
    FROM Users u
    LEFT JOIN VarificationReq vr ON vr.SubID = u.Id
    WHERE u.ID = p_ID;

    RETURN verified_status;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION usp_checkexisting(
    IN username varchar,
    IN email varchar
)
RETURNS TABLE (
    Id int,
    "Screen Name" varchar,
    "Login Name" varchar,
    Password varchar,
    Role varchar(3),
    RandOpt int,
    Varified int,
    AnonCmd bit,
    LoginDate timestamp,
    ThumbsUp int,
    VarifiedCode int
) AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM Users
    WHERE "Screen Name" = username OR "Login Name" = email;
END;
$$ LANGUAGE plpgsql;

-- FIX: Reordered parameters to place required ones first.
CREATE OR REPLACE PROCEDURE USP_BlockUser(
    IN blockee int,
    IN blocker int DEFAULT 0
)
AS $$
BEGIN
    INSERT INTO Block (BlockerId, BlockeeId)
    VALUES (blocker, blockee);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION usp_bestsender()
RETURNS TABLE (
    "Screen Name" varchar,
    ThumbsUp int
) AS $$
BEGIN
    RETURN QUERY
    SELECT u."Screen Name", u.ThumbsUp
    FROM users u
    ORDER BY u.ThumbsUp DESC
    LIMIT 5;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_AdminCommandLog()
RETURNS TABLE (
    Sender varchar,
    SubId varchar,
    GroupRefId int,
    "Decrypt" varchar,
    Content text
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COALESCE(u."screen name", 'Anon') AS Sender,
        COALESCE(su."screen name", 'Anon') AS SubId,
        c.GroupRefId,
        '' AS "Decrypt",
        c.Content
    FROM CmdLog c
    LEFT JOIN users u ON c.senderid = u.id
    LEFT JOIN users su ON c.SubId = su.id
    ORDER BY c.Id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_AddInvite(
    IN p_DomId int,
    IN p_SubName varchar(250)
)
AS $$
BEGIN
    INSERT INTO Invites (SubId, DomId)
    SELECT u.Id, p_DomId
    FROM Users u
    WHERE u."Screen Name" = p_SubName;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_AddGroup(
    IN p_user int,
    IN p_group varchar(50)
)
AS $$
DECLARE
    grouref int;
BEGIN
    SELECT RefId INTO grouref
    FROM Groups
    WHERE GroupName = p_group;

    IF (SELECT count(*) FROM UsersGroups WHERE UserId = p_user AND GroupRefId = grouref) = 0 THEN
        INSERT INTO UsersGroups (UserId, GroupRefId)
        VALUES (p_user, grouref);
    ELSE
        DELETE FROM UsersGroups
        WHERE UserId = p_user
          AND GroupRefId = grouref;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_AcceptInvite(
    IN p_SubId int,
    IN p_DomName varchar(250)
)
AS $$
DECLARE
    dom_id int;
BEGIN
    SELECT Id INTO dom_id
    FROM Users
    WHERE "Screen Name" = p_DomName;

    INSERT INTO Relationship (DomId, SubID)
    VALUES (dom_id, p_SubId);

    DELETE FROM Invites i
    USING Users u
    WHERE i.DomId = u.Id
      AND u."Screen Name" = p_DomName
      AND i.SubId = p_SubId;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetGroups(IN p_user int)
RETURNS TABLE (
    GroupName varchar,
    Allowed boolean
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        g.GroupName,
        (ug.UserId IS NOT NULL) AS Allowed
    FROM Groups g
    LEFT JOIN UsersGroups ug ON g.RefId = ug.GroupRefId AND ug.UserId = p_user;
END;
$$ LANGUAGE plpgsql;

-- FIX: Reordered parameters to place required ones first.
CREATE OR REPLACE PROCEDURE USP_Report(
    IN reportee int,
    IN reporter int DEFAULT 0,
    IN content varchar DEFAULT ''
)
AS $$
BEGIN
    CALL USP_BlockUser(reportee, reporter);

    INSERT INTO Report (ReporterId, Reportee, ReportedCommand)
    VALUES (reporter, reportee, content);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE usp_thumbsup(
    IN p_id int,
    IN p_senderid int
)
AS $$
BEGIN
    IF (SELECT count(*) FROM users WHERE id = p_senderid) > 0 THEN
        IF p_id != p_senderid THEN
            UPDATE users
            SET ThumbsUp = ThumbsUp + 1
            WHERE id = p_senderid;
        END IF;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_SendChat(
    IN p_SenderId int,
    IN p_ReceiverId int,
    IN p_message varchar,
    IN p_RecUserName varchar(250) DEFAULT ''
)
AS $$
DECLARE
    newother int;
BEGIN
    IF p_RecUserName != '' THEN
        SELECT ID INTO newother
        FROM Users
        WHERE "Screen Name" = p_RecUserName;
    ELSE
        newother := p_ReceiverId;
    END IF;

    INSERT INTO ChatLog (ReceiverID, SenderID, ChatTxt)
    VALUES (newother, p_SenderId, p_message);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_SendAppCmd(
    IN p_usernm varchar(250),
    IN p_usercmd varchar,
    IN p_all int DEFAULT 0,
    IN p_username varchar(250) DEFAULT NULL,
    IN p_password varchar(250) DEFAULT NULL,
    IN p_fileloc varchar DEFAULT NULL
)
AS $$
DECLARE
    newidins int;
    senderid int := 0;
    id int;
    numbers int;
    which int;
BEGIN
    IF p_username IS NOT NULL AND p_username != '' THEN
        SELECT subq.Id INTO senderid
        FROM (
            SELECT u.Id, u.Role
            FROM Users u
            WHERE (u."Login Name" = p_username OR u."Screen Name" = p_username)
              AND u.Password = p_password
            LIMIT 1
        ) AS subq;
    ELSIF p_username = '' THEN
        senderid := -1;
    END IF;

    IF senderid != 0 THEN
        IF p_usernm = '' THEN
            SELECT count(*) INTO numbers FROM Users WHERE RandOpt = 1;
            which := floor(random() * numbers) + 1;

            SELECT u.Id INTO id
            FROM Users u
            ORDER BY u.Id
            OFFSET which - 1
            LIMIT 1;

            INSERT INTO CommandList (Content) VALUES (p_usercmd) RETURNING CmdId INTO newidins;

            INSERT INTO ControlAppCmd (SenderId, SubId, CmdId)
            VALUES (senderid, id, newidins);

            INSERT INTO CmdLog (SenderId, SubId, Content)
            VALUES (senderid, id, p_usercmd);

        ELSIF isnumeric(p_usernm) AND CAST(p_usernm AS int) < -1 THEN
            INSERT INTO CommandList (Content) VALUES (p_usercmd) RETURNING CmdId INTO newidins;

            INSERT INTO ControlAppCmd (SenderId, SubId, CmdId, GroupRefId)
            SELECT senderid, ug.UserId, newidins, CAST(p_usernm AS int)
            FROM UsersGroups ug
            WHERE ug.GroupRefId = CAST(p_usernm AS int)
              AND ug.UserId != senderid;

            INSERT INTO CmdLog (SenderId, SubId, Content, GroupRefId)
            VALUES (senderid, senderid, p_usercmd, CAST(p_usernm AS int));
        ELSE
            SELECT sub1.Id INTO id
            FROM (
                SELECT u.Id FROM Users u WHERE u."Screen Name" = p_usernm
                UNION
                SELECT u.Id FROM Users u WHERE CAST(u.Id AS varchar(250)) = p_usernm
                LIMIT 1
            ) sub1;

            IF id IS NOT NULL THEN
                INSERT INTO CommandList (Content) VALUES (p_usercmd) RETURNING CmdId INTO newidins;

                INSERT INTO ControlAppCmd (SenderId, SubId, CmdId)
                VALUES (senderid, id, newidins);

                INSERT INTO CmdLog (SenderId, SubId, Content)
                VALUES (senderid, id, p_usercmd);
            END IF;
        END IF;
    END IF;

    DELETE FROM ControlAppCmd c
    USING CommandList cl
    WHERE c.CmdId = cl.CmdId
      AND c.SenderId < 0
      AND cl.Content LIKE '%1mp7JCqszmk=|||%';

    DELETE FROM ControlAppCmd c
    USING CommandList cl
    WHERE c.CmdId = cl.CmdId
      AND c.SenderId < 0
      AND cl.Content = '';

    DELETE FROM CommandList c
    WHERE c.CmdId NOT IN (SELECT cac.CmdId FROM ControlAppCmd cac);

END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION usp_requestvar(IN p_ID int)
RETURNS TABLE (
    "Login Name" varchar,
    VarifiedCode int
) AS $$
DECLARE
    counting int;
BEGIN
    SELECT count(*) INTO counting
    FROM VarificationReq
    WHERE SubID = p_ID;

    IF counting = 0 THEN
        INSERT INTO VarificationReq (SubID)
        VALUES (p_ID);
    END IF;

    RETURN QUERY
    SELECT u."Login Name", u.VarifiedCode
    FROM Users u
    WHERE u.ID = p_ID;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_Register(
    IN p_ScreenName varchar(255),
    IN p_Email varchar(255),
    IN p_Password varchar(255),
    IN p_Role varchar(3),
    IN p_opt int
)
AS $$
BEGIN
    INSERT INTO Users ("Screen Name", "Login Name", Password, Role, RandOpt)
    VALUES (p_ScreenName, p_Email, p_Password, p_Role, p_opt);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_Login(
    IN p_UserName varchar(255),
    IN p_Password varchar(255)
)
RETURNS TABLE (
    Id int,
    Role varchar(3),
    Varified int
) AS $$
BEGIN
    RETURN QUERY
    SELECT u.Id, u.Role, u.Varified
    FROM Users u
    WHERE u."Screen Name" = p_UserName AND u.Password = p_Password;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION Usp_GetVerificationReq()
RETURNS TABLE (
    "User" varchar,
    Email varchar,
    Code text,
    Tried int
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        u."Screen Name" AS "User",
        u."Login Name" AS Email,
        '<a href="mailto:' || u."Login Name" || '?subject=Verify &body=' || CAST(u.VarifiedCode AS varchar(4)) || '">Email</a>' AS Code,
        vr.Count AS Tried
    FROM Users u
    JOIN VarificationReq vr ON u.Id = vr.SubID
    WHERE vr.Sent = 0;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetUserSettings(IN p_userid varchar)
RETURNS SETOF Users AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM Users
    WHERE Id = CAST(p_userid AS int);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetSub(IN p_DomID int)
RETURNS TABLE (
    UserName varchar,
    SubID int
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        u."Screen Name" AS UserName,
        r.SubID
    FROM Relationship r
    JOIN Users u ON r.SubID = u.Id
    WHERE r.DomID = p_DomID;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetRels(IN p_MyID int)
RETURNS text AS $$
DECLARE
    result text;
BEGIN
    SELECT string_agg(sub.UserName, ',')
    INTO result
    FROM (
        SELECT u."Screen Name" AS UserName
        FROM Relationship r
        JOIN Users u ON r.SubID = u.Id
        WHERE r.DomID = p_MyID
        UNION
        SELECT u."Screen Name" AS UserName
        FROM Relationship r
        JOIN Users u ON r.DomID = u.Id
        WHERE r.SubID = p_MyID
    ) AS sub;

    RETURN result;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetOutstanding(IN userID int DEFAULT 0)
RETURNS TABLE (
    counting varchar(100),
    whonext varchar(250),
    varified bit,
    Thumbs int
)
AS $$
DECLARE
    Anon bit;
    vari bit;
    Thumbs_val int;
    whonext_val varchar(250);
    group_name_val varchar(50);
BEGIN
    DELETE FROM ControlAppCmd
    WHERE SenderId IN (SELECT BlockeeId FROM Block WHERE BlockerId = userID)
      AND SubId = userID;

    SELECT u.AnonCmd, u.Varified, u.ThumbsUp
    INTO Anon, vari, Thumbs_val
    FROM Users u
    WHERE u.ID = userID;

    IF Anon = B'0' THEN
        DELETE FROM ControlAppCmd
        WHERE SenderId = -1
          AND SubId = userID;
    END IF;

    SELECT COALESCE(CAST(c.GroupRefId AS varchar), CAST(c.SenderId AS varchar))
    INTO whonext_val
    FROM ControlAppCmd c
    JOIN CommandList cl ON c.CmdId = cl.CmdId
    WHERE c.SubId = userID
    ORDER BY c.SendDate
    LIMIT 1;

    UPDATE Users
    SET LoginDate = NOW()
    WHERE ID = userID;

    IF whonext_val IS NOT NULL THEN
        IF CAST(whonext_val AS int) < -1 THEN
            SELECT g.GroupName INTO group_name_val
            FROM Groups g
            WHERE g.RefId = CAST(whonext_val AS int);
            whonext_val := 'Group: ' || group_name_val;
        ELSEIF CAST(whonext_val AS int) = -1 THEN
            whonext_val := 'Anon';
        ELSE
            SELECT u."Screen Name" INTO whonext_val
            FROM users u
            JOIN Relationship r ON u.Id = r.DomID
            WHERE r.SubID = userID
            LIMIT 1;
            IF whonext_val = '' OR whonext_val IS NULL THEN
                whonext_val := 'User';
            END IF;
        END IF;
    END IF;

    RETURN QUERY
    SELECT CAST(count(*) AS varchar(100)), whonext_val, vari, Thumbs_val
    FROM ControlAppCmd
    WHERE SubId = userID;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetInvites2(IN p_SubID int)
RETURNS text AS $$
DECLARE
    result text;
BEGIN
    SELECT string_agg(sub.DomUser, ',')
    INTO result
    FROM (
        SELECT u."Screen Name" AS DomUser
        FROM Invites i
        JOIN Users u ON i.DomId = u.Id
        WHERE i.SubId = p_SubID
    ) AS sub;

    RETURN result;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION USP_GetInvites(IN p_SubID int)
RETURNS TABLE (
    DomUser varchar
) AS $$
BEGIN
    RETURN QUERY
    SELECT u."Screen Name" AS DomUser
    FROM Invites i
    JOIN Users u ON i.DomId = u.Id
    WHERE i.SubId = p_SubID;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE VIEW Vw_ControlAppCmd AS
SELECT
    c.SenderId,
    c.SubId,
    c.GroupRefId,
    cl.Content
FROM ControlAppCmd c
JOIN CommandList cl ON c.CmdId = cl.CmdId;

CREATE OR REPLACE PROCEDURE usp_veriySent(IN p_Username varchar(250))
AS $$
BEGIN
    UPDATE VarificationReq vr
    SET Sent = 1
    FROM Users u
    WHERE u.Id = vr.SubID
      AND u."Screen Name" = p_Username;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION usp_varify(
    IN p_ID int,
    IN p_varcode int
)
RETURNS text AS $$
DECLARE
    tries int;
    counting int;
BEGIN
    SELECT "Count" INTO tries
    FROM VarificationReq
    WHERE SubID = p_ID;

    IF tries <= 10 THEN
        UPDATE VarificationReq
        SET VarificationCode = p_varcode
        WHERE SubID = p_ID;

        SELECT count(*) INTO counting
        FROM Users u
        JOIN VarificationReq vr ON u.Id = vr.SubID
                                AND u.VarifiedCode = vr.VarificationCode
        WHERE u.ID = p_ID;

        IF counting = 1 THEN
            DELETE FROM VarificationReq WHERE SubID = p_ID;
            UPDATE Users SET Varified = 1 WHERE Id = p_ID;
            RETURN 'Done';
        ELSE
            UPDATE VarificationReq
            SET "Count" = "Count" + 1
            WHERE SubID = p_ID;
            RETURN 'Incorrect';
        END IF;
    ELSE
        RETURN 'Too many tried';
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION usp_userinfo(
    IN p_username varchar(255) DEFAULT '',
    IN p_id int DEFAULT 0
)
RETURNS SETOF users AS $$
BEGIN
    IF p_username != '' THEN
        RETURN QUERY SELECT * FROM users WHERE "Screen Name" = p_username;
    ELSEIF p_id != 0 THEN
        RETURN QUERY SELECT * FROM users WHERE id = p_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE USP_UpdateSettings(
    IN p_ID int,
    IN p_optin bit,
    IN p_password varchar,
    IN p_Anon bit,
    IN p_email varchar DEFAULT NULL
)
AS $$
DECLARE
    oldemail varchar;
    curvar bit;
BEGIN
    SELECT "Login Name", Varified INTO oldemail, curvar
    FROM Users
    WHERE ID = p_ID;

    IF oldemail IS DISTINCT FROM p_email THEN
        curvar := B'0';
        UPDATE Users
        SET VarifiedCode = floor(random() * 1000)
        WHERE ID = p_ID;
    END IF;

    UPDATE Users
    SET RandOpt = p_optin,
        Varified = curvar,
        Password = p_password,
        AnonCmd = p_Anon,
        "Login Name" = COALESCE(p_email, "Login Name")
    WHERE ID = p_ID;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE usp_tidyup()
AS $$
BEGIN
    DELETE FROM CommandList
    WHERE senddate < NOW() - INTERVAL '5 day';

    DELETE FROM ControlAppCmd
    WHERE subid IN (
        SELECT cac.subid
        FROM ControlAppCmd cac
        LEFT JOIN CommandList cl ON cac.CmdId = cl.CmdId
        WHERE cl.CmdId IS NULL
    );

    DELETE FROM CommandList
    WHERE CmdId NOT IN (SELECT cac.CmdId FROM ControlAppCmd cac);

    DELETE FROM VarificationReq
    WHERE VerifyDate < NOW() - INTERVAL '5 day';

    DELETE FROM usersgroups
    WHERE grouprefid = -11 AND userid != 1;
END;
$$ LANGUAGE plpgsql;