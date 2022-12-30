autoload -U add-zsh-hook

hist_record() {
    __HIST_COMMAND_INDEX=0
    hist record "$1"
}
add-zsh-hook preexec hist_record

hist-backward-widget() {
    __HIST_COMMAND_INDEX=$((__HIST_COMMAND_INDEX+1))
    command=$(hist get --index $__HIST_COMMAND_INDEX)
    ret=$?
    if [ $ret -ne 0 ]; then
        zle reset-prompt
        return $ret
    fi
    BUFFER=$command
    zle vi-add-eol
}

zle -N                      hist-backward-widget
bindkey -M emacs "^[[A"     hist-backward-widget
bindkey -M viins "^[[A"     hist-backward-widget
bindkey -M vicmd "^[[A"     hist-backward-widget

hist-forward-widget() {
    if [ $__HIST_COMMAND_INDEX -gt 0 ]; then
        __HIST_COMMAND_INDEX=$((__HIST_COMMAND_INDEX-1))
        command=$(hist get --index $__HIST_COMMAND_INDEX)
        ret=$?
        if [ $ret -ne 0 ]; then
            zle reset-prompt
            return $ret
        fi
        BUFFER=$command
        zle vi-add-eol
    fi
}

zle -N                      hist-forward-widget
bindkey -M emacs "^[[B"     hist-forward-widget
bindkey -M viins "^[[B"     hist-forward-widget
bindkey -M vicmd "^[[B"     hist-forward-widget

TRAPINT() {
    __HIST_COMMAND_INDEX=0
    # return 128 plus the signal number
    return $(( 128 + $1 ))
}

hist-search-widget() {
    zle push-line
    command=$(hist search)
    ret=$?
    if [ $ret -ne 0 ]; then
        zle reset-prompt
        return $ret
    fi
    BUFFER=$command
    zle vi-add-eol
}

zle -N                  hist-search-widget
bindkey -M emacs '^H'   hist-search-widget
bindkey -M viins '^H'   hist-search-widget
bindkey -M vicmd '^H'   hist-search-widget
