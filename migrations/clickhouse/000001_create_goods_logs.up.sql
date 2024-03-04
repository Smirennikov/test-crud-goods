CREATE TABLE IF NOT EXISTS goods_logs (
    Id Int32,
    ProjectId Int32,
    Name String,
    Description String,
    Priority Int32,
    Removed UInt8,
    EventTime DateTime,

    INDEX idx_id Id TYPE minmax GRANULARITY 1,
    INDEX idx_project_id ProjectId TYPE minmax GRANULARITY 1,
    INDEX idx_name Name TYPE bloom_filter GRANULARITY 1
) ENGINE = MergeTree()
ORDER BY (EventTime)
