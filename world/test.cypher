// MATCH (n)
// DETACH DELETE n;

// // ------------------------------------------------------------------

// MATCH (p:Player {username:'test@test.com'})
// OPTIONAL MATCH (p)-[r:IN]->(s)
// MATCH (d:Location {name:'kitchens'})
// DELETE r
// CREATE (n)-[:IN]->(d);

// MATCH (p:Player)
// RETURN p;

// MATCH (:Player {username:'becolem@microsoft.com'})-[:IN]->(l:Location)
// MATCH (l)-[v]->(e:Location)
// RETURN v,e

MATCH (:Player {username:'becolem@microsoft.com'})-[:IN]->(l:Location) MATCH (l)-[r]->(:Location) RETURN r