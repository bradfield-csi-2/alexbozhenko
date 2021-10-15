
/* delete existing tables to make this script stateless.
`set rebuild false` below to avoid rebuilding */
create or replace function rebuild() returns void as $$
declare
    num_employees int;
    num_bonuses int;
    min_salary int;
    max_salary int;
begin
    num_employees := 500000;
    num_bonuses := 2000000;
    min_salary := 70000;
    max_salary := 80000;

    raise notice 'Rebuilding test data';
    set client_min_messages=warning;
    drop table if exists department cascade;
    drop table if exists employee cascade;
    drop table if exists bonus;


    /* (re)create tables */
    create table if not exists department (
        dep_id serial primary key,
        name varchar(100)
    );
    create table if not exists employee (
        emp_id serial primary key,
        dep_id integer references department ( dep_id ),
        manager_id integer,
        name varchar(100),
        salary integer
    );
    create table if not exists bonus (
        bonus_id serial primary key,
        emp_id integer references employee ( emp_id ),
        amount integer,
        time timestamp
    );

    /* create an index on employee salary */
    create index employee_manager_salary on employee (manager_id);
    --create index bonus_emp_id_time on bonus (emp_id, time);
    create index bonus_emp_id on bonus (emp_id);
    create index emp_dep_id on employee(dep_id);
--    create index employee_salary on employee (salary);
--    create index employee_salary_dep on employee (salary, dep_id);
--    create index employee_manager_salary on employee (manager_id, salary);

    /* we don't need many departments to keep this interesting */
    insert into department ( name ) values
        ( 'sales' ),
        ( 'marketing' ),
        ( 'engineering' );


    /* create a bunch of employees, mostly in sales */
    insert into employee ( dep_id, manager_id, name, salary ) (
        select
            ('{1,1,1,2,3}'::int[])[floor(random()*5) + 1],  /* skew towards sales */
            1 + n % 10,  /* first 10 employees manage all the others */
            md5(random()::text),  /* name is just a made up string */
            min_salary + random() * (max_salary - min_salary)  /* uniform dist of salaries in range */
        from generate_series(1, num_employees) as n
    );


    /* Add some bonus payments */
    insert into bonus ( emp_id, amount, time ) (
        select
            1 + floor(random() * num_employees),
            10000 * random(),  /* bonus amount is uniformly in [0, 10000) */
            now() - random() * (now()+'720 days' - now())  /* all paid in last 30 days */
        from generate_series(1, num_bonuses)
    );


    /* ensure that the catalog is up to date */
    analyze;
end
$$ language plpgsql;


/* rebuild test data: comment this out to use existing data set */
\t on
--select rebuild();
\t off

/* average employee salary... note what happens when we have 3-4x the number of employees */
-- explain select count(*), avg(salary) from employee;

-- 1. What is the total bonus amount for each employee?

/*
explain analyze select employee.emp_id, sum(amount) 
from employee left join bonus 
on employee.emp_id = bonus.emp_id
group by employee.emp_id;
*/

-- 2. Which employees received the highest total compensation last year, and how much was it?

/* explain analyze select e.emp_id, e.name, e.salary + sum(b.amount) as tc 
from employee e
join bonus b on e.emp_id = b.emp_id
where b.time between '01-01-2020'::timestamp and '31-12-2020'::timestamp
group by e.emp_id -- turns out, starting with pg 9.1, having pkey in "group by" is enough
order by tc desc
limit 10; */

--3. What is the total bonus amount awarded by department?

explain analyze verbose select d.name, sum(amount) as total_bonus 
from department d
left join employee e 
    on e.dep_id = d.dep_id
left join bonus b 
    on b.emp_id = e.emp_id
group by d.dep_id;

--4. Which employees earned more than their managers?


/* top salaries in sales... note what indexes would help here? */
--explain select * from employee where employee.dep_id = 1 order by employee.salary limit 10;

/* how many employees earn more than their managers? what indexes would help here? */
--explain select e.name, m.name
--from employee as e, employee as m
--where e.manager_id = m.emp_id
--and e.salary > m.salary;

/* top employee comps including bonus... note what happens when we have 3-4x the employees */
-- explain select employee.emp_id, employee.salary + sum(bonus.amount) as total_comp
-- from employee, bonus
-- where bonus.emp_id = employee.emp_id
-- group by employee.emp_id, employee.salary
-- order by total_comp
-- limit 10;

/* see some stats */
/*
select histogram_bounds
from pg_stats
where tablename='employee' and attname='salary';
*/
