# Go Calendar


# WIP

* Task repetition logic, how do the task repeat and persist through an application?
    * How can I identify which task is the origin and which are replicated?
        - Compound:ID. Make unique ID identifier for tasks.
        - task_id-<serial-number>-o / task_id-<serial-number>-c
        - only save tasks completed to the db, by default uncompleted without draft state will be generated on runtime.
        - Task's global state will be cached. 
    * How can I mutate all of them? How can I distiguish them even when they are the same?


# Calendar Sequence Diagrams

## Create Task

```mermaid
sequenceDiagram
    actor U as User
    participant C as Calendar
    participant M as Month
    participant D as Day
    participant T as Task


    U->>C: Call Calendar.AddTask method
    activate C
    activate C
    activate C
    activate C
    #activate C
    #C->>T: Call Task.NewTask method
    #activate T
    #activate T
    #alt Is valid Task
    #    T-->>C: Return Task
    #    deactivate T
    #else Is invalid
    #    T-->>C: Throw error
    #    deactivate T
    #    C-->>U: Print error
    #    deactivate C
    #end
    
    #C->>C: Find Task Month
    #C->>M: Month[N] AddTaskToDay(day int, task *Task)
    par Obtain tasks time.Time through its channel (datesChan)
        C->>T: Call Task.SearchRepetition() <-chan time.Time
        activate T
        T->>T: Declare variables for isRepeating and the interval time
        T->>T: Declare a timeChan for results
        alt Is task repeating
            T->>T: Declare variable for lastDay for validation
            T->>T: Spin a goroutine
            Note over T,T: Spin a goroutine which declares an init currenTime from the Task time.Time
            Note over T,T: Validates if the currentTime isn't the last day or after it
            Note over T,T: Adds the currentTime to the timeChan 
            Note over T,T: Sets the currentTime to the new time added with the interval
        else Is not repeating
            T->>T: Sends own time.Time and closes the channel
        end
        T-->>C: Send datesChan
        deactivate T
    and Obtain a channel of Tasks created with the values sent to the datesChan
        C->>T: Call Task.CreateTasksFromDates(datesChan <-chan time.Time) <-chan *Task
        activate T
        T->>T: Declare a taskChan for results
        T->>T: Initialize taskIDFactory (copy)
        T->>T: Spin a goroutine
        Note over T,T: Creates a taskID for the task to be created
        Note over T,T: Creates a task and obtain the error
        Note over T,T: Omit error - for now
        Note over T,T: Send task through channel        
        T-->>C: Send tasksChan 
        deactivate T
    end

    #activate M
    
    loop Every item in tasksChan
        C->>C: Call Calendar.addMonth(month time.Month) (*Month, error)
        alt Is no error
            C->>C: Assign month local variable
        else Is error
            C-->>U: Throw error
        end
        C->>M: Call Month.addTaskToDay(day int, task *Task) error
        activate M
        M-->>C: Return error
        deactivate M

        alt addTaskToDay: Is no error
            C->>C: Continue loop
            Note over C, M: Task is assigned
        else addTaskToDay: Is error
            alt addTaskToDay: Is ErrDayNotFound
                C->>M: Call month.addDay(day int) (*Day, error)
                activate M
                M-->>C: Return Day and Error
                deactivate M

                alt addDay: Is error
                    C-->>U: Return error
                deactivate C
                else addDay: Is no error
                    C->>D: Call Day.addTask(task *Task) error
                    activate D
                    D-->>C: Return error
                    
                    deactivate D
                    alt Is error
                        C-->>U: Return error
                        deactivate C
                    else Is no error
                        C-->>C: Continue loop
                    end
                end
            else addTaskToDay: Is other error
                C-->>U: Return error
                deactivate C
            end
  
        end
        C-->>C: Finishes loop
    end

    #deactivate M

    C-->>U: Return nil 
    deactivate C   
```