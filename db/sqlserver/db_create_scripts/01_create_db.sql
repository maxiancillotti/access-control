USE [master]
GO
/****** Object:  Database [access_control]    Script Date: 7/26/2022 22:36:58 ******/
CREATE DATABASE [access_control]
-- CONTAINMENT = NONE
-- ON  PRIMARY 
--( NAME = N'Auth', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL14.SQLEXPRESS\MSSQL\DATA\Auth.mdf' , SIZE = 8192KB , MAXSIZE = UNLIMITED, FILEGROWTH = 65536KB )
-- LOG ON 
--( NAME = N'Auth_log', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL14.SQLEXPRESS\MSSQL\DATA\Auth_log.ldf' , SIZE = 8192KB , MAXSIZE = 2048GB , FILEGROWTH = 65536KB )
GO
ALTER DATABASE [access_control] SET COMPATIBILITY_LEVEL = 140
GO
IF (1 = FULLTEXTSERVICEPROPERTY('IsFullTextInstalled'))
begin
EXEC [access_control].[dbo].[sp_fulltext_database] @action = 'enable'
end
GO
ALTER DATABASE [access_control] SET ANSI_NULL_DEFAULT OFF 
GO
ALTER DATABASE [access_control] SET ANSI_NULLS OFF 
GO
ALTER DATABASE [access_control] SET ANSI_PADDING OFF 
GO
ALTER DATABASE [access_control] SET ANSI_WARNINGS OFF 
GO
ALTER DATABASE [access_control] SET ARITHABORT OFF 
GO
ALTER DATABASE [access_control] SET AUTO_CLOSE OFF 
GO
ALTER DATABASE [access_control] SET AUTO_SHRINK OFF 
GO
ALTER DATABASE [access_control] SET AUTO_UPDATE_STATISTICS ON 
GO
ALTER DATABASE [access_control] SET CURSOR_CLOSE_ON_COMMIT OFF 
GO
ALTER DATABASE [access_control] SET CURSOR_DEFAULT  GLOBAL 
GO
ALTER DATABASE [access_control] SET CONCAT_NULL_YIELDS_NULL OFF 
GO
ALTER DATABASE [access_control] SET NUMERIC_ROUNDABORT OFF 
GO
ALTER DATABASE [access_control] SET QUOTED_IDENTIFIER OFF 
GO
ALTER DATABASE [access_control] SET RECURSIVE_TRIGGERS OFF 
GO
ALTER DATABASE [access_control] SET  DISABLE_BROKER 
GO
ALTER DATABASE [access_control] SET AUTO_UPDATE_STATISTICS_ASYNC OFF 
GO
ALTER DATABASE [access_control] SET DATE_CORRELATION_OPTIMIZATION OFF 
GO
ALTER DATABASE [access_control] SET TRUSTWORTHY OFF 
GO
ALTER DATABASE [access_control] SET ALLOW_SNAPSHOT_ISOLATION OFF 
GO
ALTER DATABASE [access_control] SET PARAMETERIZATION SIMPLE 
GO
ALTER DATABASE [access_control] SET READ_COMMITTED_SNAPSHOT OFF 
GO
ALTER DATABASE [access_control] SET HONOR_BROKER_PRIORITY OFF 
GO
ALTER DATABASE [access_control] SET RECOVERY SIMPLE 
GO
ALTER DATABASE [access_control] SET  MULTI_USER 
GO
ALTER DATABASE [access_control] SET PAGE_VERIFY CHECKSUM  
GO
ALTER DATABASE [access_control] SET DB_CHAINING OFF 
GO
ALTER DATABASE [access_control] SET FILESTREAM( NON_TRANSACTED_ACCESS = OFF ) 
GO
ALTER DATABASE [access_control] SET TARGET_RECOVERY_TIME = 60 SECONDS 
GO
ALTER DATABASE [access_control] SET DELAYED_DURABILITY = DISABLED 
GO
ALTER DATABASE [access_control] SET QUERY_STORE = OFF
GO
USE [access_control]
GO
/****** Object:  Table [dbo].[Admins]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Admins](
	[ID] [int] IDENTITY(1,1) NOT NULL,
	[Username] [nvarchar](100) NOT NULL,
	[PasswordHash] [nvarchar](64) NOT NULL,
	[PasswordSalt] [nvarchar](64) NOT NULL,
	[Enabled] [bit] NOT NULL,
	[Deleted] [bit] NOT NULL,
 CONSTRAINT [PK_Admins_] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[HttpMethods]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[HttpMethods](
	[ID] [int] IDENTITY(1,1) NOT NULL,
	[Name] [nvarchar](50) NOT NULL,
 CONSTRAINT [PK_Permissions] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY],
 CONSTRAINT [IX_Unique_Name] UNIQUE NONCLUSTERED 
(
	[Name] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Resources]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Resources](
	[ID] [int] IDENTITY(1,1) NOT NULL,
	[Path] [nvarchar](200) NOT NULL,
	[Deleted] [bit] NOT NULL,
 CONSTRAINT [PK_Resources] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[Users]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Users](
	[ID] [int] IDENTITY(1,1) NOT NULL,
	[Username] [nvarchar](100) NOT NULL,
	[PasswordHash] [nvarchar](64) NOT NULL,
	[PasswordSalt] [nvarchar](64) NOT NULL,
	[Enabled] [bit] NOT NULL,
	[Deleted] [bit] NOT NULL,
 CONSTRAINT [PK_Users_] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[UsersRESTPermissions]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[UsersRESTPermissions](
	[UserID] [int] NOT NULL,
	[ResourceID] [int] NOT NULL,
	[HttpMethodID] [int] NOT NULL,
 CONSTRAINT [PK_Users_Resources_Methods_1] PRIMARY KEY CLUSTERED 
(
	[UserID] ASC,
	[ResourceID] ASC,
	[HttpMethodID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
SET ANSI_PADDING ON
GO
/****** Object:  Index [IX_Admins_Username_Unique]    Script Date: 7/26/2022 22:36:58 ******/
CREATE UNIQUE NONCLUSTERED INDEX [IX_Admins_Username_Unique] ON [dbo].[Admins]
(
	[Username] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, DROP_EXISTING = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
SET ANSI_PADDING ON
GO
/****** Object:  Index [IX_Unique_Path]    Script Date: 7/26/2022 22:36:58 ******/
CREATE NONCLUSTERED INDEX [IX_Unique_Path] ON [dbo].[Resources]
(
	[Path] ASC,
	[Deleted] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, DROP_EXISTING = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
SET ANSI_PADDING ON
GO
/****** Object:  Index [IX_Users_Username_Unique]    Script Date: 7/26/2022 22:36:58 ******/
CREATE UNIQUE NONCLUSTERED INDEX [IX_Users_Username_Unique] ON [dbo].[Users]
(
	[Username] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, DROP_EXISTING = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Admins] ADD  CONSTRAINT [DF_Admins_Deleted]  DEFAULT ((0)) FOR [Deleted]
GO
ALTER TABLE [dbo].[Resources] ADD  CONSTRAINT [DF_Resources_Deleted]  DEFAULT ((0)) FOR [Deleted]
GO
ALTER TABLE [dbo].[Users] ADD  CONSTRAINT [DF_Users_Deleted]  DEFAULT ((0)) FOR [Deleted]
GO
ALTER TABLE [dbo].[UsersRESTPermissions]  WITH CHECK ADD  CONSTRAINT [FK_UsersRESTPermissions_HttpMethods] FOREIGN KEY([HttpMethodID])
REFERENCES [dbo].[HttpMethods] ([ID])
GO
ALTER TABLE [dbo].[UsersRESTPermissions] CHECK CONSTRAINT [FK_UsersRESTPermissions_HttpMethods]
GO
ALTER TABLE [dbo].[UsersRESTPermissions]  WITH CHECK ADD  CONSTRAINT [FK_UsersRESTPermissions_Resources] FOREIGN KEY([ResourceID])
REFERENCES [dbo].[Resources] ([ID])
GO
ALTER TABLE [dbo].[UsersRESTPermissions] CHECK CONSTRAINT [FK_UsersRESTPermissions_Resources]
GO
ALTER TABLE [dbo].[UsersRESTPermissions]  WITH CHECK ADD  CONSTRAINT [FK_UsersRESTPermissions_Users] FOREIGN KEY([UserID])
REFERENCES [dbo].[Users] ([ID])
GO
ALTER TABLE [dbo].[UsersRESTPermissions] CHECK CONSTRAINT [FK_UsersRESTPermissions_Users]
GO
/****** Object:  StoredProcedure [dbo].[AdminsDelete]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[AdminsDelete]
	@id int
AS
BEGIN

	UPDATE Admins
	
	SET	Deleted = 1

	WHERE ID = @id

END
GO
/****** Object:  StoredProcedure [dbo].[AdminsExists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[AdminsExists]
	@id int,
	@exists bit out
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF EXISTS (
	
		SELECT *
		FROM Admins
		WHERE ID = @id
			and Deleted = 0
	)

	--Return values supported by Go's driver in an inconvenient way because it adds a dependency I want to avoid
	--https://github.com/denisenkom/go-mssqldb#return-status
	/*
		RETURN 1
	ELSE
		RETURN 0
	*/
		SET @exists = 1
	ELSE
		SET @exists = 0
	
END
GO
/****** Object:  StoredProcedure [dbo].[AdminsInsert]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[AdminsInsert]
	@username nvarchar(100),
	@passwordHash nvarchar(64),
	@passwordSalt nvarchar(64),
	@enabled bit,
	@ID int output
AS
BEGIN

	INSERT INTO Admins
	(
		Username,
		PasswordHash,
		PasswordSalt,
		Enabled
	)
	VALUES
	(
		@username,
		@passwordHash,
		@passwordSalt,
		@enabled
	)

	---------------------------------

	SET @ID = SCOPE_IDENTITY()

END
GO
/****** Object:  StoredProcedure [dbo].[AdminsSelect]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.30
-- =============================================
CREATE PROCEDURE [dbo].[AdminsSelect]
	@id int
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		ID,
		Username,
		PasswordHash,
		PasswordSalt,
		[Enabled]
	FROM Admins
	WHERE ID = @id
		and Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[AdminsSelectByUsername]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2021.10.24
-- =============================================
CREATE PROCEDURE [dbo].[AdminsSelectByUsername]
	@username nvarchar(100)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		ID,
		Username,
		PasswordHash,
		PasswordSalt,
		[Enabled]
	FROM Admins
	WHERE Username = @username
		and Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[AdminsUpdateEnableState]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[AdminsUpdateEnableState]
	@id int,
	@enabled bit
AS
BEGIN

	UPDATE Admins
	
	SET	Enabled = @enabled
		
	WHERE ID = @id

END
GO
/****** Object:  StoredProcedure [dbo].[AdminsUpdatePassword]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[AdminsUpdatePassword]
	@id int,
	@passwordHash nvarchar(64),
	@passwordSalt nvarchar(64)
AS
BEGIN

	UPDATE Admins
	
	SET	PasswordHash = @passwordHash,
		PasswordSalt = @passwordSalt
	
	WHERE ID = @id

END
GO
/****** Object:  StoredProcedure [dbo].[AdminsUsernameExists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[AdminsUsernameExists]
	@username nvarchar(100),
	@exists bit out
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF EXISTS (
	
		SELECT *
		FROM Admins
		WHERE Username = @username
			--don't check for deleted status
			--we need to know if it ever existed, because username cannot be duplicated.
	)

	--Return values supported by Go's driver in an inconvenient way because it adds a dependency I want to avoid
	--https://github.com/denisenkom/go-mssqldb#return-status
	/*
		RETURN 1
	ELSE
		RETURN 0
	*/
		SET @exists = 1
	ELSE
		SET @exists = 0
	
END
GO
/****** Object:  StoredProcedure [dbo].[HttpMethodsSelectAll]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.30
-- =============================================
CREATE PROCEDURE [dbo].[HttpMethodsSelectAll]
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT *
	FROM HttpMethods
	ORDER BY ID
END
GO
/****** Object:  StoredProcedure [dbo].[HttpMethodsSelectByName]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2021.10.24
-- =============================================
CREATE PROCEDURE [dbo].[HttpMethodsSelectByName]
	@name nvarchar(50)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT *
	FROM HttpMethods
	WHERE [Name] = @name
END
GO
/****** Object:  StoredProcedure [dbo].[ResourcesDelete]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[ResourcesDelete]
	@id int
AS
BEGIN

	UPDATE [Resources]
	
	SET	Deleted = 1

	WHERE ID = @id

END
GO
/****** Object:  StoredProcedure [dbo].[ResourcesExists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[ResourcesExists]
	@id int,
	@exists bit out
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF EXISTS (
	
		SELECT *
		FROM Resources
		WHERE ID = @id
			and Deleted = 0
	)

	--Return values supported by Go's driver in an inconvenient way because it adds a dependency I want to avoid
	--https://github.com/denisenkom/go-mssqldb#return-status
	/*
		RETURN 1
	ELSE
		RETURN 0
	*/
		SET @exists = 1
	ELSE
		SET @exists = 0
	
END
GO
/****** Object:  StoredProcedure [dbo].[ResourcesInsert]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[ResourcesInsert]
	@path nvarchar(200),
	@ID int output
AS
BEGIN

	INSERT INTO [Resources]
	(
		[Path]
	)
	VALUES
	(
		@path
	)

	---------------------------------

	SET @ID = SCOPE_IDENTITY()

END
GO
/****** Object:  StoredProcedure [dbo].[ResourcesPathExists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[ResourcesPathExists]
	@path nvarchar(200),
	@exists bit out
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF EXISTS (
	
		SELECT *
		FROM Resources
		WHERE [Path] = @path
			--don't check for deleted status
			--we need to know if it ever existed, because resource name cannot be duplicated.
	)

	--Return values supported by Go's driver in an inconvenient way because it adds a dependency I want to avoid
	--https://github.com/denisenkom/go-mssqldb#return-status
	/*
		RETURN 1
	ELSE
		RETURN 0
	*/
		SET @exists = 1
	ELSE
		SET @exists = 0
	
END
GO
/****** Object:  StoredProcedure [dbo].[ResourcesSelect]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[ResourcesSelect]
	@id int
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		ID,
		[Path]
	FROM [Resources]
	WHERE ID = @id
		and Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[ResourcesSelectAll]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.30
-- =============================================
CREATE PROCEDURE [dbo].[ResourcesSelectAll]
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		ID,
		[Path]
	FROM Resources
	WHERE Deleted = 0
	ORDER BY ID
END
GO
/****** Object:  StoredProcedure [dbo].[ResourcesSelectByPath]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[ResourcesSelectByPath]
	@path nvarchar(200)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		ID,
		[Path]
	FROM [Resources]
	WHERE Path = @path
		and Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[UsersDelete]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[UsersDelete]
	@id int
AS
BEGIN

	UPDATE Users
	
	SET	Deleted = 1

	WHERE ID = @id

END
GO
/****** Object:  StoredProcedure [dbo].[UsersExists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersExists]
	@id int,
	@exists bit out
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF EXISTS (
	
		SELECT *
		FROM Users
		WHERE ID = @id
			and Deleted = 0
	)

	--Return values supported by Go's driver in an inconvenient way because it adds a dependency I want to avoid
	--https://github.com/denisenkom/go-mssqldb#return-status
	/*
		RETURN 1
	ELSE
		RETURN 0
	*/
		SET @exists = 1
	ELSE
		SET @exists = 0
	
END
GO
/****** Object:  StoredProcedure [dbo].[UsersInsert]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[UsersInsert]
	@username nvarchar(100),
	@passwordHash nvarchar(64),
	@passwordSalt nvarchar(64),
	@enabled bit,
	@ID int output
AS
BEGIN

	INSERT INTO Users
	(
		Username,
		PasswordHash,
		PasswordSalt,
		Enabled
	)
	VALUES
	(
		@username,
		@passwordHash,
		@passwordSalt,
		@enabled
	)

	---------------------------------

	SET @ID = SCOPE_IDENTITY()

END
GO
/****** Object:  StoredProcedure [dbo].[UsersRESTPermissions_Delete]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersRESTPermissions_Delete]
	@userID int,
	@resourceID int,
	@methodID int
AS
BEGIN

	DELETE
	FROM [UsersRESTPermissions]
	
	WHERE
		[UserID] = @userID and
		[ResourceID] = @resourceID and
		[HttpMethodID] = @methodID
END
GO
/****** Object:  StoredProcedure [dbo].[UsersRESTPermissions_Exists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersRESTPermissions_Exists]
	@userID int,
	@resourceID int,
	@methodID int,
	@exists bit out
AS

BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF EXISTS (
	
		SELECT *

		FROM [UsersRESTPermissions]
	
		WHERE
			[UserID] = @userID and
			[ResourceID] = @resourceID and
			[HttpMethodID] = @methodID
	)

	--Return values supported by Go's driver in an inconvenient way because it adds a dependency I want to avoid
	--https://github.com/denisenkom/go-mssqldb#return-status
	/*
		RETURN 1
	ELSE
		RETURN 0
	*/
		SET @exists = 1
	ELSE
		SET @exists = 0
	
END
GO
/****** Object:  StoredProcedure [dbo].[UsersRESTPermissions_Insert]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersRESTPermissions_Insert]
	@userID int,
	@resourceID int,
	@methodID int
AS
BEGIN

	INSERT INTO [UsersRESTPermissions]
	(
		[UserID],
		[ResourceID],
		[HttpMethodID]
	)
	VALUES
	(
		@userID,
		@resourceID,
		@methodID
	)

END
GO
/****** Object:  StoredProcedure [dbo].[UsersRESTPermissions_RelationshipsExists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersRESTPermissions_RelationshipsExists]
	@userID int,
	@resourceID int,
	@methodID int,
	@exists int out
AS

BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF NOT (
		EXISTS (
			SELECT *
			FROM [Users]
			WHERE [ID] = @userID
		) --or @userID = 0
	)
	BEGIN
		SET @exists = -1
		RETURN
	END

	IF NOT (
		EXISTS (
			SELECT *
			FROM [Resources]
			WHERE [ID] = @resourceID
		) --or @resourceID = 0
	)
	BEGIN
		SET @exists = -2
		RETURN
	END

	IF NOT (
		EXISTS (
			SELECT *
			FROM [HttpMethods]
			WHERE [ID] = @methodID
		)
		--or @methodID = 0
	)
	BEGIN
		SET @exists = -3
		RETURN
	END

	SET @exists = 1
END
GO
/****** Object:  StoredProcedure [dbo].[UsersRESTPermissions_SelectAllByUserID]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersRESTPermissions_SelectAllByUserID]
	@userID int
AS
BEGIN

	SELECT
		ResourceID,
		HttpMethodID
	
	FROM [UsersRESTPermissions]

	WHERE UserID = @userID

	-- IMPORTANT: The service layer depends on this order to traverse the resulting data structure correctly and efficiently.
	ORDER BY UserID, ResourceID, HttpMethodID
	

END
GO
/****** Object:  StoredProcedure [dbo].[UsersRESTPermissions_SelectAllPathsMethodsByUserID]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.01.11
-- =============================================
CREATE PROCEDURE [dbo].[UsersRESTPermissions_SelectAllPathsMethodsByUserID]
	@userID int
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		R.[Path] 'Resource',
		M.[Name] 'Method'

	FROM UsersRESTPermissions URP
		inner join Users U on U.ID = URP.UserID
		inner join Resources R on R.ID = URP.ResourceID
		inner join HttpMethods M on M.ID = URP.HttpMethodID

	WHERE U.ID = @userID
		and U.Deleted = 0
		and R.Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[UsersRESTPermissions_SelectAllWithDescriptionsByUserID]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersRESTPermissions_SelectAllWithDescriptionsByUserID]
	@userID int
AS
BEGIN

	SELECT

		URP.ResourceID,
		R.[Path] 'ResourcePath',
		URP.HttpMethodID,
		M.[Name] 'HttpMethodName'
	
	FROM [UsersRESTPermissions] URP
		inner join Users U on U.ID = URP.UserID
		inner join Resources R on R.ID = URP.ResourceID
		inner join HttpMethods M on M.ID = URP.HttpMethodID

	WHERE UserID = @userID

	-- IMPORTANT: The service layer depends on this order to traverse the resulting data structure correctly and efficiently.
	ORDER BY UserID, ResourceID, HttpMethodID
	

END
GO
/****** Object:  StoredProcedure [dbo].[UsersSelect]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.30
-- =============================================
CREATE PROCEDURE [dbo].[UsersSelect]
	@id int
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		ID,
		Username,
		PasswordHash,
		PasswordSalt,
		[Enabled]
	FROM Users
	WHERE ID = @id
		and Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[UsersSelectAllResourcesPermissionsByUserID]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.01.11
-- =============================================
CREATE PROCEDURE [dbo].[UsersSelectAllResourcesPermissionsByUserID]
	@userID int
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT
		R.[Name] 'Resource',
		P.[Name] 'Permission'

	FROM Users_Resources_Permissions URP
		inner join Users U on U.ID = URP.UserID
		inner join Resources R on R.ID = URP.ResourceID
		inner join [Permissions] P on P.ID = URP.PermissionID

	WHERE U.ID = @userID
		and U.Deleted = 0
		and R.Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[UsersSelectByUsername]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2021.10.24
-- =============================================
CREATE PROCEDURE [dbo].[UsersSelectByUsername]
	@username nvarchar(100)
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	SELECT 
		ID,
		Username,
		PasswordHash,
		PasswordSalt,
		[Enabled]
	FROM Users
	WHERE Username = @username
		and Deleted = 0
END
GO
/****** Object:  StoredProcedure [dbo].[UsersUpdateEnableState]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[UsersUpdateEnableState]
	@id int,
	@enabled bit
AS
BEGIN

	UPDATE Users
	
	SET	Enabled = @enabled
		
	WHERE ID = @id

END
GO
/****** Object:  StoredProcedure [dbo].[UsersUpdatePassword]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.28
-- =============================================
CREATE PROCEDURE [dbo].[UsersUpdatePassword]
	@id int,
	@passwordHash nvarchar(64),
	@passwordSalt nvarchar(64)
AS
BEGIN

	UPDATE Users
	
	SET	PasswordHash = @passwordHash,
		PasswordSalt = @passwordSalt
	
	WHERE ID = @id

END
GO
/****** Object:  StoredProcedure [dbo].[UsersUsernameExists]    Script Date: 7/26/2022 22:36:58 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
-- =============================================
-- Author:		Maximiliano Ancillotti
-- Create date: 2022.03.29
-- =============================================
CREATE PROCEDURE [dbo].[UsersUsernameExists]
	@username nvarchar(100),
	@exists bit out
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON;

	IF EXISTS (
	
		SELECT *
		FROM Users
		WHERE Username = @username
			--don't check for deleted status
			--we need to know if it ever existed, because username cannot be duplicated.
	)

	--Return values supported by Go's driver in an inconvenient way because it adds a dependency I want to avoid
	--https://github.com/denisenkom/go-mssqldb#return-status
	/*
		RETURN 1
	ELSE
		RETURN 0
	*/
		SET @exists = 1
	ELSE
		SET @exists = 0
	
END
GO
USE [master]
GO
ALTER DATABASE [access_control] SET  READ_WRITE 
GO
