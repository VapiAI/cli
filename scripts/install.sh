#!/usr/bin/env bash
set -euo pipefail

platform=$(uname -ms)

if [[ ${OS:-} = Windows_NT ]]; then
  if [[ $platform != MINGW64* ]]; then
    powershell -c "irm vapi.ai/install.ps1|iex"
    exit $?
  fi
fi

# Reset
Color_Off=''

# Regular Colors
Red=''
Green=''
Dim=''

# Bold
Bold_White=''
Bold_Green=''

if [[ -t 1 ]]; then
    # Reset
    Color_Off='\033[0m'

    # Regular Colors
    Red='\033[0;31m'
    Green='\033[0;32m'
    Dim='\033[0;2m'

    # Bold
    Bold_Green='\033[1;32m'
    Bold_White='\033[1m'
fi

error() {
    echo -e "${Red}error${Color_Off}:" "$@" >&2
    exit 1
}

info() {
    echo -e "${Dim}$@ ${Color_Off}"
}

info_bold() {
    echo -e "${Bold_White}$@ ${Color_Off}"
}

success() {
    echo -e "${Green}$@ ${Color_Off}"
}

command -v curl >/dev/null ||
    error 'curl is required to install vapi'

command -v tar >/dev/null ||
    error 'tar is required to install vapi'

case $platform in
'Darwin x86_64')
    target=Darwin_x86_64
    ;;
'Darwin arm64')
    target=Darwin_arm64
    ;;
'Linux aarch64' | 'Linux arm64')
    target=Linux_arm64
    ;;
'MINGW64'*)
    target=Windows_x86_64
    ;;
'Linux x86_64' | *)
    target=Linux_x86_64
    ;;
esac

GITHUB=${GITHUB-"https://github.com"}
github_repo="$GITHUB/VapiAI/cli"

# Get latest version
if [[ $# = 0 ]]; then
    version=$(curl -s "https://api.github.com/repos/VapiAI/cli/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [[ -z $version ]]; then
        error "Failed to fetch latest version"
    fi
    vapi_uri="$github_repo/releases/download/$version/cli_$target.tar.gz"
else
    vapi_uri="$github_repo/releases/download/$1/cli_$target.tar.gz"
fi

install_env=VAPI_INSTALL
bin_env=\$$install_env/bin

install_dir=${!install_env:-$HOME/.vapi}
bin_dir=$install_dir/bin
exe=$bin_dir/vapi

if [[ ! -d $bin_dir ]]; then
    mkdir -p "$bin_dir" ||
        error "Failed to create install directory \"$bin_dir\""
fi

curl --fail --location --progress-bar --output "$exe.tar.gz" "$vapi_uri" ||
    error "Failed to download vapi from \"$vapi_uri\""

tar -xzf "$exe.tar.gz" -C "$bin_dir" ||
    error 'Failed to extract vapi'

# Handle different possible binary names
if [[ -f "$bin_dir/vapi.exe" ]]; then
    mv "$bin_dir/vapi.exe" "$exe"
elif [[ -f "$bin_dir/vapi" ]]; then
    mv "$bin_dir/vapi" "$exe"
else
    error 'Failed to find vapi binary in extracted files'
fi

chmod +x "$exe" ||
    error 'Failed to set permissions on vapi executable'

rm "$exe.tar.gz"

tildify() {
    if [[ $1 = $HOME/* ]]; then
        local replacement=\~/
        echo "${1/$HOME\//$replacement}"
    else
        echo "$1"
    fi
}

success "vapi was installed successfully to $Bold_Green$(tildify "$exe")"

if command -v vapi >/dev/null; then
    echo "Run 'vapi --help' to get started"
    exit
fi

refresh_command=''

tilde_bin_dir=$(tildify "$bin_dir")
quoted_install_dir=\"${install_dir//\"/\\\"}\"

if [[ $quoted_install_dir = \"$HOME/* ]]; then
    quoted_install_dir=${quoted_install_dir/$HOME\//\$HOME/}
fi

echo

case $(basename "$SHELL") in
fish)
    commands=(
        "set --export $install_env $quoted_install_dir"
        "set --export PATH $bin_env \$PATH"
    )

    fish_config=$HOME/.config/fish/config.fish
    tilde_fish_config=$(tildify "$fish_config")

    if [[ -w $fish_config ]]; then
        {
            echo -e '\n# vapi'
            for command in "${commands[@]}"; do
                echo "$command"
            done
        } >>"$fish_config"

        info "Added \"$tilde_bin_dir\" to \$PATH in \"$tilde_fish_config\""
        refresh_command="source $tilde_fish_config"
    else
        echo "Manually add the directory to $tilde_fish_config (or similar):"
        for command in "${commands[@]}"; do
            info_bold "  $command"
        done
    fi
    ;;
zsh)
    commands=(
        "export $install_env=$quoted_install_dir"
        "export PATH=\"$bin_env:\$PATH\""
    )

    zsh_config=$HOME/.zshrc
    tilde_zsh_config=$(tildify "$zsh_config")

    if [[ -w $zsh_config ]]; then
        {
            echo -e '\n# vapi'
            for command in "${commands[@]}"; do
                echo "$command"
            done
        } >>"$zsh_config"

        info "Added \"$tilde_bin_dir\" to \$PATH in \"$tilde_zsh_config\""
        refresh_command="exec $SHELL"
    else
        echo "Manually add the directory to $tilde_zsh_config (or similar):"
        for command in "${commands[@]}"; do
            info_bold "  $command"
        done
    fi
    ;;
bash)
    commands=(
        "export $install_env=$quoted_install_dir"
        "export PATH=\"$bin_env:\$PATH\""
    )

    bash_configs=(
        "$HOME/.bashrc"
        "$HOME/.bash_profile"
    )

    if [[ ${XDG_CONFIG_HOME:-} ]]; then
        bash_configs+=(
            "$XDG_CONFIG_HOME/.bash_profile"
            "$XDG_CONFIG_HOME/.bashrc"
            "$XDG_CONFIG_HOME/bash_profile"
            "$XDG_CONFIG_HOME/bashrc"
        )
    fi

    set_manually=true
    for bash_config in "${bash_configs[@]}"; do
        tilde_bash_config=$(tildify "$bash_config")

        if [[ -w $bash_config ]]; then
            {
                echo -e '\n# vapi'
                for command in "${commands[@]}"; do
                    echo "$command"
                done
            } >>"$bash_config"

            info "Added \"$tilde_bin_dir\" to \$PATH in \"$tilde_bash_config\""
            refresh_command="source $bash_config"
            set_manually=false
            break
        fi
    done

    if [[ $set_manually = true ]]; then
        echo "Manually add the directory to $tilde_bash_config (or similar):"
        for command in "${commands[@]}"; do
            info_bold "  $command"
        done
    fi
    ;;
*)
    echo 'Manually add the directory to ~/.bashrc (or similar):'
    info_bold "  export $install_env=$quoted_install_dir"
    info_bold "  export PATH=\"$bin_env:\$PATH\""
    ;;
esac

echo
info "To get started, run:"
echo

if [[ $refresh_command ]]; then
    info_bold "  $refresh_command"
fi

info_bold "  vapi --help"
