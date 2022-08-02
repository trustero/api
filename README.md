# API
An API to interface with the Trustero platform

# Diagrams

## Discovered
```mermaid
sequenceDiagram

    actor alice
    participant receptor_sdk.Receptor
    participant receptor_v1.ReceptorService
    participant agent.UpdaterService
    
    
    alice ->> receptor_sdk.Receptor: Receptor.Report(reports)
    receptor_sdk.Receptor ->> receptor_v1.ReceptorService: ReceptorService.Report(finding)
    receptor_v1.ReceptorService ->> agent.UpdaterService: UpdaterService.Report(Discovery)
    agent.UpdaterService ->> receptor_v1.ReceptorService: DiscoveryId
    receptor_v1.ReceptorService ->> receptor_sdk.Receptor: StringValue
    receptor_sdk.Receptor ->> alice: done
```
### Verified
```mermaid
sequenceDiagram

    actor alice
    participant receptor_sdk.Receptor
    participant receptor_v1.ReceptorService
    participant receptor.VerifyCredentials
    
    
    alice ->> receptor_sdk.Receptor: Receptor.Report(reports)
    activate receptor_sdk.Receptor
    receptor_sdk.Receptor ->> receptor_v1.ReceptorService: ReceptorService.Verified(Credentials)
    receptor_v1.ReceptorService ->> receptor.VerifyCredentials: Credentials
    receptor.VerifyCredentials ->> receptor_v1.ReceptorService: Credentials
    receptor_v1.ReceptorService ->> receptor_sdk.Receptor: StringValue
    receptor_sdk.Receptor ->> alice: done
```

## Report
```mermaid
sequenceDiagram

    actor alice
    participant receptor_sdk.Receptor
    participant receptor_v1.ReceptorService
    participant agent.UpdaterService
    participant attachment.AttachmentService
    
    
    alice ->> receptor_sdk.Receptor: Receptor.Report(reports)
    activate receptor_sdk.Receptor
    receptor_sdk.Receptor ->> receptor_v1.ReceptorService: ReceptorService.Report(finding)
    activate receptor_v1.ReceptorService
    receptor_v1.ReceptorService ->> agent.UpdaterService: UpdaterService.Report(Discovery)
    agent.UpdaterService ->> receptor_v1.ReceptorService: DiscoveryId
    receptor_v1.ReceptorService ->> attachment.AttachmentService: AttachmentService.AddEvidene(evidence,discoveryId)
    attachment.AttachmentService ->>receptor_v1.ReceptorService: Document
    receptor_v1.ReceptorService ->> receptor_sdk.Receptor: StringValue
    receptor_sdk.Receptor ->> alice: done
```


### Appendix A: Legacy implementation of Scan

```mermaid
sequenceDiagram

    actor alice
    participant cmd
    participant Receptor
       
    alice ->> cmd: %> scan <token>
    cmd ->> Receptor: Scan(credentials)
    Receptor ->> cmd: []*agent.Services
    cmd ->> agent.UpdaterService.Report(): []*agent.Services
    agent.UpdaterService.Report() ->> cmd: DiscoveryId
    cmd ->> Receptor: FindEvidence([]*agent.Services)
    Receptor ->> cmd: []*evidence.Evidence
    cmd ->> attachment.AttachmentService.AddEvidence(): (evidence,discoveryId)
    attachment.AttachmentService.AddEvidence() ->>cmd: Document
    cmd ->> alice: done
```

## Message Passing example: Scan
```mermaid
sequenceDiagram

    actor alice
    participant main
    participant cmd
    participant ReceptorClient
    participant ReceptorServer

    alice->>main: "scan <token>"
    main->>cmd: cmd.Executer(receptorImpl)"
    cmd->>cmd: "set receptorImpl"
    cmd->>cmd: "call cmd.scan"
    cmd->>cmd: "receptorImpl.Verify(credentials)"
    cmd->>cmd: "ok, error"
    cmd->>cmd: "new receptor_v1.Credentials(...)"
    cmd->>ReceptorClient: "receptor.Verified(credentials)"
    ReceptorClient->>ReceptorServer: "receptor.Verified(credentials)"
    ReceptorServer->>ntrced/receptor/receptor.go: "ReceptorServer.Verified(credentials)"
    ntrced/receptor/receptor.go->>ReceptorServer: "&empty.Empty{}, nil"
    ReceptorServer->>ReceptorClient: "&empty.Empty{}, nil"
    ReceptorClient->>cmd: "&empty.Empty{}, nil"
    cmd->>ReceptorClient: "receptor..Notify scan success"
    ReceptorClient->>ReceptorServer: "scan success"
    ReceptorServer->>ReceptorClient: "&empty.Empty{}, nil"
```