autoload -U add-zsh-hook

function TRAPINT() {
    __HIST_COMMAND_INDEX=0
    [[ -o zle ]] && zle reset-prompt
    # return 128 plus the signal number
    return $(( 128 + $1 ))
}

function hist-record() {
    __HIST_COMMAND_INDEX=0
    hist record "$1"
}
add-zsh-hook preexec hist-record

hist-backward-widget() {
    if [ ${__HIST_DIRECTION:="backward"} = "forward" ]; then
        __HIST_DIRECTION="backward"
        __HIST_COMMAND_INDEX=$((__HIST_COMMAND_INDEX+1))
    fi
    command=$(hist get --index ${__HIST_COMMAND_INDEX:=0})
    ret=$?
    if [ $ret -ne 0 ]; then
        zle reset-prompt
        return $ret
    fi
    __HIST_COMMAND_INDEX=$((__HIST_COMMAND_INDEX+1))
    BUFFER=$command
    zle end-of-line
}

zle -N                      hist-backward-widget
bindkey -M emacs "^[[A"     hist-backward-widget
bindkey -M viins "^[[A"     hist-backward-widget
bindkey -M vicmd "^[[A"     hist-backward-widget

hist-forward-widget() {
    if [ ${__HIST_DIRECTION:="forward"} = "backward" ]; then
        __HIST_COMMAND_INDEX=$((__HIST_COMMAND_INDEX-1))
        __HIST_DIRECTION="forward"
    fi
    if [ ${__HIST_COMMAND_INDEX:=0} -gt 0 ]; then
        __HIST_COMMAND_INDEX=$((__HIST_COMMAND_INDEX-1))
        command=$(hist get --index $__HIST_COMMAND_INDEX)
        ret=$?
        if [ $ret -ne 0 ]; then
            zle reset-prompt
            return $ret
        fi
        BUFFER=$command
        zle end-of-line
    else
        __HIST_COMMAND_INDEX=0
        zle kill-whole-line
    fi
}

zle -N                      hist-forward-widget
bindkey -M emacs "^[[B"     hist-forward-widget
bindkey -M viins "^[[B"     hist-forward-widget
bindkey -M vicmd "^[[B"     hist-forward-widget

hist-search-widget() {
    zle reset-prompt
    command=$(hist search)
    ret=$?
    if [ $ret -ne 0 ]; then
        zle reset-prompt
        return $ret
    fi
    BUFFER=$command
    zle end-of-line
}

zle -N                  hist-search-widget
bindkey -M emacs '^R'   hist-search-widget
bindkey -M viins '^R'   hist-search-widget
bindkey -M vicmd '^R'   hist-search-widget
