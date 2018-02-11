# rbt-control

Tool to quickly encapsulates somo common aws-operations for quick automation

## Usage

```
docker run \
    -e AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
    atarax/rbl-control:intermediate \
    /rbl-control -r "eu-central-1" -c {{COMMAND}}
```

Command can be either list, create or destory. 

- create creates one instance
- destroy terminates all instances created by this tool
- list lists all instances
