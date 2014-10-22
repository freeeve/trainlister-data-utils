LOAD CSV FROM "file:///Users/wfreeman/graphhack/temp.csv" as coll
WITH toInt(coll[0]) as id, coll[1] as direction, coll[2] as type, coll[3] as saturday, coll[4] as name, toFloat(coll[5]) as mile, toInt(coll[6]) as hours, toInt(coll[7]) as minutes
MERGE (t:Train {id:id, direction:direction, type:type, saturday:case when saturday = "false" then null else true end})
MERGE (s:Stop {name:name, mile:mile})
MERGE (t)-[:STOPS_AT {hours:hours, minutes:minuts}]->(s);
