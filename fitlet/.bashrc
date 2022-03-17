# If not running interactively, don't do anything
case $- in
    *i*) ;;
      *) return;;
esac

# don't put duplicate lines or lines starting with space in the history.
HISTCONTROL=ignoreboth

# append to the history file, don't overwrite it
shopt -s histappend

# for setting history length see HISTSIZE and HISTFILESIZE in bash(1)
HISTSIZE=100000
HISTFILESIZE=200000

# ignore certain commands with history
HISTIGNORE='ls:history:exit'

# timestamp on history
HISTTIMEFORMAT='%F %T '

# check the window size after each command and, if necessary,
# update the values of LINES and COLUMNS.
shopt -s checkwinsize

# make less more friendly for non-text input files, see lesspipe(1)
[ -x /usr/bin/lesspipe ] && eval "$(SHELL=/bin/sh lesspipe)"

# set variable identifying the chroot you work in (used in the prompt below)
if [ -z "${debian_chroot:-}" ] && [ -r /etc/debian_chroot ]; then
    debian_chroot=$(cat /etc/debian_chroot)
fi

# set a fancy prompt (non-color, unless we know we "want" color)
case "$TERM" in
    xterm-color|*-256color) color_prompt=yes;;
esac

# makes a colored prompt
force_color_prompt=yes

if [ -n "$force_color_prompt" ]; then
    if [ -x /usr/bin/tput ] && tput setaf 1 >&/dev/null; then
	# We have color support; assume it's compliant with Ecma-48
	# (ISO/IEC-6429). (Lack of such support is extremely rare, and such
	# a case would tend to support setf rather than setaf.)
	color_prompt=yes
    else
	color_prompt=
    fi
fi

###
### Color and prompt
###

# changes colorprofile based on user
#1;32m is normal green with no backing
PROMPT_NORMAL="1;32m" # default green with no back
PROMPT_RED="1;41m"
PROMPT_GREEN="1;42m"

PROMPT_FLUKE="1;30;43m" # yellow bg with bold black text

if [ "$(whoami)" = "fitlet" ]; then
      usercolor=$PROMPT_GREEN

elif [ "$(whoami)" = "root" ]; then
      usercolor=$PROMPT_RED
else
      usercolor=$PROMPT_NORMAL
fi

# sets color profile
if [ "$color_prompt" = yes ]; then
    PS1='${debian_chroot:+($debian_chroot)}\[\e[${usercolor}\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ '
else
    PS1='${debian_chroot:+($debian_chroot)}\u@\h:\w\$ '
fi
unset color_prompt force_color_prompt

# colored GCC warnings and errors
export GCC_COLORS='error=01;31:warning=01;35:note=01;36:caret=01;32:locus=01:quote=01'

###
### aliases
###

# enable color support of ls and also add handy aliases
if [ -x /usr/bin/dircolors ]; then
    test -r ~/.dircolors && eval "$(dircolors -b ~/.dircolors)" || eval "$(dircolors -b)"
    alias ls='ls --color=auto'
    alias dir='dir --color=auto'
    alias vdir='vdir --color=auto'
    alias grep='grep --no-messages --color=auto'
    alias fgrep='fgrep --no-messages --color=auto'
    alias egrep='egrep --no-messages --color=auto'
fi

# some more aliases
alias ll='ls -alhF'
alias la='ls -A'
alias l='ls -CF'

# enable programmable completion features
if ! shopt -oq posix; then
  if [ -f /usr/share/bash-completion/bash_completion ]; then
    . /usr/share/bash-completion/bash_completion
  elif [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
  fi
fi

# make vim default text editor
export EDITOR=/usr/bin/vim

# displays fortune at each login
/usr/games/fortune -s | /usr/games/cowsay

# alias for easy updates/upgrades
alias upgrade='sudo apt-get update && sudo apt-get upgrade'

# vi mode
set -o vi

# path variables
export PATH=$PATH:/usr/local/go/bin