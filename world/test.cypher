// MATCH (n)
// DETACH DELETE n;

// // ------------------------------------------------------------------

// MATCH (n:Player {username:'test@test.com'})
// OPTIONAL MATCH (n)-[:IN]->(s)
// MATCH (d:Location {name:'kitchens'})
// DELETE r
// CREATE (n)-[:IN]->(d);

// MATCH (p:Player)
// RETURN p;

MATCH (l:Player) RETURN l;

// MATCH (p:Player)-->(l:Location)
// RETURN p, l;