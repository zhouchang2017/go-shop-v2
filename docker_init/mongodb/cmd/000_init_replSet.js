rs.initiate({
        _id: "rs0",
        members: [
            {_id: 0, host: "mongo_rs1"},
            {_id: 1, host: "mongo_rs2"},
            {_id: 2, host: "mongo_rs3", arbiterOnly: true}
        ]
    }
);
