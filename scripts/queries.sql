-- 1. Invested by brokerage

select
  	i.brokerage,
	count(*) as counter,
  	(
  		sum(i.total_value) - (
  		  select COALESCE(sum(ivt.total_value), 0) as total
  		  from investments ivt 
          where ivt.operation_type = 'sell' 
  				and ivt.brokerage = i.brokerage
  				and (ivt.due_date is null or ivt.due_date = '' or ivt.due_date >= DATE('now'))
  		  )
  	) as total_invested
from investments i
where 
  i.operation_type = 'buy'
  and (due_date is null or due_date = '' or due_date >= DATE('now'))
group by 1;