# GCP publisher (with GIN endpoint)

Esportare la chiave pubblica

     export GOOGLE_APPLICATION_CREDENTIALS=...

Per eseguire

     ./web training-gcp-309207 topic-frank

Per pubblicare

     curl -X POST  http://localhost:8080/content \
         -d '{"id":"23846582137", "key":"key1", "user":"frank", "abstract":"an abstract", "content":"some interesting content", "filename":"content-x.json"}'

Pubblicazione bulk

     curl -X POST  http://localhost:8080/content/bulk \
        -d '[
            {"id":"0000001", "key":"key1", "user":"frank", "abstract":"an abstract 1", "content":"some interesting content part 1", "filename":"content-frank-1.json"},
            {"id":"0000002", "key":"key1", "user":"frank", "abstract":"an abstract 2", "content":"some interesting content part 2", "filename":"content-frank-2.json"},
            {"id":"0000003", "key":"key2", "user":"mike", "abstract":"shiny launder", "content":"foo bar boo", "filename":"content-mike-2.json"},
            {"id":"0000004", "key":"key3", "user":"john", "abstract":"mistic pizza", "content":"foo bar boo", "filename":"content-john-2.json"},
            {"id":"0000005", "key":"key1", "user":"frank", "abstract":"an abstract 3", "content":"some interesting content part 3", "filename":"content-frank-3.json"}
            ]'


