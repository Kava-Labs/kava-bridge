# States


## Per Signing Session State

States per signing session. There may be multiple sessions running in parallel.

```mermaid
stateDiagram-v2
    state if_leader_state <<choice>>
    state "Is leader?" as IsLeader

    [*] --> IsLeader: Start signing session
    IsLeader --> if_leader_state: pick leader
    if_leader_state --> LeaderWaitingForCandidatesState: is leader
    if_leader_state --> CandidateWaitingForLeaderState: not leader
    LeaderWaitingForCandidatesState --> LeaderWaitingForCandidatesState: AddCandidateEvent

    state if_participant_state <<choice>>
    state "Is participant?" as IsParticipant

    LeaderWaitingForCandidatesState --> SigningState: t + 1 picked
    CandidateWaitingForLeaderState --> IsParticipant: StartSignerEvent
    IsParticipant --> if_participant_state
    if_participant_state --> SigningState: true
    if_participant_state --> DoneWithoutSignatureState: false

    SigningState --> SigningState: AddSigningPartEvent
    SigningState --> DoneWithSignatureState: Signing done
    SigningState --> ErrorState: Signing error
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
