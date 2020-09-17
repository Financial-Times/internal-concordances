# UPP Internal Concordances API

Simplified Concordances API with extra Concept data for Internal use.

## Primary URL
https://api.ft.com/internalconcordances

## Service Tier
Bronze

## Lifecycle Stage
Production

## Delivered By
content

## Supported By
content

## Known About By

- elitsa.pavlova
- kalin.arsov
- miroslav.gatsanoga
- ivan.nikolov
- marina.chompalova
- donislav.belev
- mihail.mihaylov
- boyko.boykov
- dimitar.terziev

## Host Platform
AWS

## Architecture
Helm chart deployed in Kubernetes. Аccepts a list of concept IDs, concords them, and returns them as an array of UPP concepts.

## Contains Personal Data
Yes

## Contains Sensitive Data
No

## Dependencies
* up-csa
* public-concordances-api

## Failover Architecture Type
ActiveActive

## Failover Process Type
FullyAutomated

## Failback Process Type
FullyAutomated

## Data Recovery Process Type

NotApplicable

## Data Recovery Details

The service does not store data, so it does not require any data recovery steps.

## Release Process Type

PartiallyAutomated

## Rollback Process Type
Manual

## Key Management Process Type

Manual

## Key Management Details

To access the service clients need to provide basic auth credentials.
To rotate credentials you need to login to a particular cluster and update varnish-auth secrets.

## Monitoring
<table class="table table-bordered" style="border-color: rgb(238, 238, 238); margin: 1em 0px; color: rgb(34, 34, 34); font-family: HelveticaNeue-Light, &quot;Helvetica Neue Light&quot;, &quot;Helvetica Neue&quot;, Arial, sans-serif; font-size: 15px; background-color: rgb(255, 255, 255);"><tbody><tr><td style="line-height: 1.5; padding: 10px 6px; border-bottom-color: rgb(238, 238, 238);">Cluster Health Check</td><td style="line-height: 1.5; padding: 10px 6px; border-bottom-color: rgb(238, 238, 238);"><p><a href="https://upp-prod-delivery-eu.ft.com/__health" target="_blank" style="color: rgb(38, 116, 122);">Prod-EU</a></p><p><a href="https://upp-prod-delivery-us.ft.com/__health" target="_blank" style="color: rgb(38, 116, 122);">Prod-US</a></p></td></tr><tr><td style="line-height: 1.5; padding: 10px 6px; border-bottom-color: rgb(238, 238, 238);">Service Health Check</td><td style="line-height: 1.5; padding: 10px 6px; border-bottom-color: rgb(238, 238, 238);"><p>Cluster aggregate healthchecks are located here:</p><p><a href="https://upp-prod-delivery-eu.ft.com/__health" target="_blank">Prod-EU</a></p><p><a href="https://upp-prod-delivery-us.ft.com/__health" target="_blank">Prod-US</a></p><p>Key "internal-concordances" into the search bar, and click the links to find individual system healthchecks&nbsp;&nbsp;&nbsp;&nbsp;</p></td></tr></tbody></table><br>

## First Line Troubleshooting
<div style="">
                                                                                                                                                                                                                                                                                                                                                                           <p style="color: rgb(34, 34, 34); font-family: HelveticaNeue-Light, &quot;Helvetica Neue Light&quot;, &quot;Helvetica Neue&quot;, Arial, sans-serif; font-size: 15px;">If the /__health endpoint is reporting issues with the service's Technical Dependencies, please refer to the Dewey guides for those services:</p><ul style=""><li style=""><font color="#222222" face="HelveticaNeue-Light, Helvetica Neue Light, Helvetica Neue, Arial, sans-serif"><span style="font-size: 15px;"><a href="https://dewey.in.ft.com/view/system/up-csa" target="_blank">Concept Search API&nbsp;</a></span></font></li><li style=""><font color="#222222" face="HelveticaNeue-Light, Helvetica Neue Light, Helvetica Neue, Arial, sans-serif"><span style="font-size: 15px;"><a href="https://dewey.in.ft.com/view/system/public-concordances-api" target="_blank">Public Concordances API</a>&nbsp;</span></font></li></ul><p style="color: rgb(34, 34, 34); font-family: HelveticaNeue-Light, &quot;Helvetica Neue Light&quot;, &quot;Helvetica Neue&quot;, Arial, sans-serif; font-size: 15px;"><br></p>