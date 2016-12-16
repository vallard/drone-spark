The Cisco Spark plugin posts build status messages to your Spark Room.  The below pipeline configuration demonstrates simple usage: 

```
pipeline:
	spark:
		image: vallard/drone-spark
		room: "CI Builds"
		
```

The room should be by name, not ID. 

Example configuration with a custom message template:

```
pipeline:
	spark:
		image: vallard/drone-spark
		room: "OurBuildRoom"
		template: |
			{{ #success build.status }}
				build {{ build.number }} succeeded.  Good job.
			{{ else }}
				build {{ build.number }} failed.  Fix me please.
			{{ /success }}
			
```
Spark also supports markdown.  For custom markdown templates you would use the following: 

```
pipeline
	spark:
		image: vallard/drone-spark
		room: OurBuildRoom
		markdown: |
			{{ #success build.status }}
				## Hurray!\n __{{build.number}}__ succeeded!.  Good job.
			{{ else }}
				## OH NO!!!!!\n __{{ build.number}}__ failed.  **{{ build.author }}** Please fix your build!
			{{ /success }}
```

## Secrets
The Spark plugin supports reading credentials from the Drone secret store.  This is strongly recommended instead of storing credentials in the pipeline configuration in plain text. 

```diff
pipeline:
	spark:
		image: vallard/drone-spark
		room: MyBuildRoom
-      	token: Y2jse.... 		
```

You should instead set your spark token using the drone secrets with: 

```
drone secret add --image=vallard/drone-spark myorg/myproject SPARK_TOKEN Ym4jZ...
```
Then be sure to sign your repo before committing: 

```
drone sign myorg/myproject
```

Drone spark will then use the spark token to post messages to your room. 

## Parameter Reference
room: The room name 

token:  The spark token.  Don't use this, use the secret instead. 

template: overwrite the default message template

markdown: overwrite the default message with a markdown template. 

## Template Reference

This code was incorporated from the [plugins/slack](https://github.com/drone-plugins/drone-slack) repo. 

uppercasefirst
: converts the first letter of a string to uppercase

uppercase
: converts a string to uppercase

lowercase
: converts a string to lowercase. Example `{{lowercase build.author}}`

datetime
: converts a unix timestamp to a date time string. Example `{{datetime build.started}}`

success
: returns true if the build is successful

failure
: returns true if the build is failed

truncate
: returns a truncated string to n characters. Example `{{truncate build.sha 8}}`

urlencode
: returns a url encoded string

since
: returns a duration string between now and the given timestamp. Example `{{since build.started}}`
