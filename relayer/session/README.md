# States


## Per Signing Session State

States per signing session. There may be multiple sessions running in parallel.

```mermaid
stateDiagram-v2
    state if_leader_state <<choice>>
    [*] --> IsLeader
    IsLeader --> if_leader_state
    if_leader_state --> LeaderWaitingForCandidatesState: is leader
    if_leader_state --> CandidatesWaitingForLeaderState: not leader
    LeaderWaitingForCandidatesState --> LeaderWaitingForCandidatesState: AddCandidateEvent
    LeaderWaitingForCandidatesState --> Signing: t + 1 picked
    CandidatesWaitingForLeaderState --> Signing: StartSignerEvent
    Signing --> Signing: AddSigningPartEvent
    Signing --> Done: Signing output
```

### 1. PickLeader

1. All peers determine the leader
2. Send `JoinSessionMessage` to leader
3. Receive `SigningPartyStartMessage`
4. If in session, go to `Signing` State
5. If not in session, end state

### 2. Signing

1. Send and receive `SigningPartMessage` until signature output
2. Signature output go to `Broadcast` state
3. Send `SigningOutputMessage` to check all peers have the same signature
4. Go to `Broadcast`

### 3. Broadcast

1. Broadcast message to destination chain Kava / Ethereum.
