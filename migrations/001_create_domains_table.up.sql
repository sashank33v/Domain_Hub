Create Table domains(
    id serial primary key,
    domain_name varchar(200) not null unique,
    registrar varchar(100) not null,
    expiry_date date not null,
    status varchar(20) not null
);
