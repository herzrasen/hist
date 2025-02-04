function hist_reset
    set -e hist_counter
    commandline -r ""
end

function hist_up
    if not set -q hist_counter
        set -g hist_counter 0
    else
        set -g hist_counter (math $hist_counter + 1)
    end
    commandline -r (hist get --index $hist_counter)
end

function hist_down
    if not set -q hist_counter
        set -g hist_counter 0
    end
    if test $hist_counter -gt 0
        set -g hist_counter (math $hist_counter - 1)
        commandline -r (hist get --index $hist_counter)
    else
        commandline -r ""
    end
end

function hist_enter
    set current_commandline (commandline -b)
    if test -n "$current_commandline"
        hist record "$current_commandline"
    end
    set -e hist_counter
    commandline -f execute
end

function hist_search
    set current_commandline (commandline -b)
    set search_result (hist search $current_commandline)
    if test -n "$search_result"
        commandline -r $search_result
    end
end

bind \e\[A 'hist_up'
bind \e\[B 'hist_down'
bind \r 'hist_enter'
bind \cc 'hist_reset'
bind \cr 'hist_search'
