USE [access_control]
GO

insert into Admins (Username, PasswordHash, PasswordSalt, [Enabled], [Deleted])
values ('admin', '$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC', 'sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S', 1, 0)