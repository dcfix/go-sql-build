use internal
go

if exists (select name from sysobjects where name = 'stp_bogus')
    drop proc stp_bogus
go

create proc stp_bogus
as
select fname, lname from employee where lname = 'fix'
go


