rs.initiate({
        _id: "rs0",
        members: [
            {_id: 0, host: "mongodb-primary:30000"},
            {_id: 1, host: "mongodb-secondary:30001"},
            {_id: 2, host: "mongodb-arbiter:30002", arbiterOnly: true}
        ]
    }
);
