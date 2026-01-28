drop schema if exists app_code;
create schema app_code;

{{ template "funcs/tg_updated_at.sql" . }}
