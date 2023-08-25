--[[
lvim is the global options object

Linters should be
filled in as strings with either
a global executable or a path to
an executable
]]
-- THESE ARE EXAMPLE CONFIGS FEEL FREE TO CHANGE TO WHATEVER YOU WANT

-- general
lvim.log.level = "warn"
lvim.format_on_save.enabled = true
lvim.colorscheme = "lunar"
-- to disable icons and use a minimalist setup, uncomment the following
-- lvim.use_icons = false

-- keymappings [view all the defaults by pressing <leader>Lk]
lvim.leader = "space"
lvim.keys.normal_mode["<C-s>"] = ":w<cr>"
lvim.keys.normal_mode["<Esc><Esc>"] = "<C-\\><C-n>"
lvim.keys.normal_mode["<S-l>"] = ":BufferLineCycleNext<CR>"
lvim.keys.normal_mode["<S-h>"] = ":BufferLineCyclePrev<CR>"

-- Remaps for the refactoring operations currently offered by the plugin 'refactoring'
lvim.keys.visual_block_mode["<leader>rb"] = "<Cmd>lua require('refactoring').refactor('Extract Block')<CR>"
lvim.keys.normal_mode["<leader>rbf"] = "<Cmd>lua require('refactoring').refactor('Extract Block To File')<CR>"
lvim.keys.normal_mode["<leader>ri"] = "<Cmd>lua require('refactoring').refactor('Inline Variable')<CR>"
lvim.keys.normal_mode["<leader>rr"] = "<Cmd>lua require('telescope').extensions.refactoring.refactors()<CR>"

-- Remaps for the refactoring operations currently offered by the plugin 'refactoring'
lvim.keys.visual_block_mode["<leader>re"] = "<Esc><Cmd>lua require('refactoring').refactor('Extract Function')<CR>"
lvim.keys.visual_block_mode["<leader>rf"] =
"<Esc><Cmd>lua require('refactoring').refactor('Extract Function To File')<CR>"
lvim.keys.visual_block_mode["<leader>rv"] = "<Esc><Cmd>lua require('refactoring').refactor('Extract Variable')<CR>"
lvim.keys.visual_block_mode["<leader>ri"] = "<Esc><Cmd>lua require('refactoring').refactor('Inline Variable')<CR>"
lvim.keys.visual_mode["<leader>re"] = "<Esc><Cmd>lua require('refactoring').refactor('Extract Function')<CR>"
lvim.keys.visual_mode["<leader>rf"] = "<Esc><Cmd>lua require('refactoring').refactor('Extract Function To File')<CR>"
lvim.keys.visual_mode["<leader>rv"] = "<Esc><Cmd>lua require('refactoring').refactor('Extract Variable')<CR>"
lvim.keys.visual_mode["<leader>ri"] = "<Esc><Cmd>lua require('refactoring').refactor('Inline Variable')<CR>"
-- Always paste the text which was yanked originally.
lvim.keys.visual_mode["p"] = '"0p'
lvim.keys.visual_mode["P"] = '"0p'
lvim.keys.visual_block_mode["p"] = '"0p'
lvim.keys.visual_block_mode["P"] = '"0p'
-- Spectre global search and replace
lvim.keys.normal_mode["<leader>S"] = "<Cmd>lua require('spectre').open()<CR>"
-- Search for the current word
lvim.keys.normal_mode["<leader>sw"] = "<Cmd>lua require('spectre').open_visual({select_word=true})<CR>"
-- Search for word in current file
lvim.keys.normal_mode["<leader>s"] = "<Esc>:lua require('spectre').open_visual()<CR>"
-- Search for word in current file
lvim.keys.normal_mode["<leader>sp"] = "viw:lua require('spectre').open_file_search()<CR>"



-- nnoremap <leader>S <cmd>lua require('spectre').open()<CR>

-- "search current word
-- nnoremap <leader>sw <cmd>lua require('spectre').open_visual({select_word=true})<CR>
-- vnoremap <leader>s <esc>:lua require('spectre').open_visual()<CR>
-- "  search in current file
-- nnoremap <leader>sp viw:lua require('spectre').open_file_search()<cr>


-- -- Rename variable under cursor
lvim.keys.normal_mode["<leader>rn"] = "<cmd>lua vim.lsp.buf.rename()<CR>"

-- unmap a default keymapping
-- vim.keymap.del("n", "<C-Up>")
-- override a default keymapping
-- lvim.keys.normal_mode["<C-q>"] = ":q<cr>" -- or vim.keymap.set("n", "<C-q>", ":q<cr>" )

-- Change Telescope navigation to use j and k for navigation and n and p for history in both input and normal mode.
-- we use protected-mode (pcall) just in case the plugin wasn't loaded yet.
local _, actions = pcall(require, "telescope.actions")
lvim.builtin.telescope.defaults.mappings = {
  -- for input mode
  i = {
    ["<C-j>"] = actions.move_selection_next,
    ["<C-k>"] = actions.move_selection_previous,
    ["<C-n>"] = actions.cycle_history_next,
    ["<C-p>"] = actions.cycle_history_prev,
  },
  -- for normal mode
  n = {
    ["<C-j>"] = actions.move_selection_next,
    ["<C-k>"] = actions.move_selection_previous,
  },
}

-- Change theme settings
-- lvim.builtin.theme.options.dim_inactive = true
-- lvim.builtin.theme.options.style = "storm"

-- Use which-key to add extra bindings with the leader-key prefix
lvim.builtin.which_key.mappings["P"] = { "<cmd>Telescope projects<CR>", "Projects" }
lvim.builtin.which_key.mappings["t"] = {
  name = "+Trouble",
  r = { "<cmd>Trouble lsp_references<cr>", "References" },
  f = { "<cmd>Trouble lsp_definitions<cr>", "Definitions" },
  d = { "<cmd>Trouble document_diagnostics<cr>", "Diagnostics" },
  q = { "<cmd>Trouble quickfix<cr>", "QuickFix" },
  l = { "<cmd>Trouble loclist<cr>", "LocationList" },
  w = { "<cmd>Trouble workspace_diagnostics<cr>", "Workspace Diagnostics" },
}

-- Help for refactoring
-- Note that these options are available in NORMAL, VISUAL and VISUAL_BLOCK mode as well
lvim.builtin.which_key.mappings["r"] = {
  name = "+Refactor",
  b = {
    "<Cmd>lua require('refactoring').refactor('Extract Block')<CR>",
    "Extract Block To File",
    f = { "<Cmd>lua require('refactoring').refactor('Extract Block To File')<CR>", "Extract Block To File" }
  },
  i = { "<Cmd>lua require('refactoring').refactor('Inline Variable')<CR>", "Inline Variable" },
  r = { "<Cmd>lua require('telescope').extensions.refactoring.refactors()<CR>", "Show Refactoring Options" },
  e = { "<Esc><Cmd>lua require('refactoring').refactor('Extract Function')<CR>", "Extract Function" },
  f = { "<Esc><Cmd>lua require('refactoring').refactor('Extract Function To File')<CR>", "Extract Function To File" },
  v = { "<leader>re [[ <Esc><Cmd>lua require('refactoring').refactor('Extract Function')<CR>]]", "Extract Function" },
}

-- lvim.builtin.which_key.mappings["C"] = {
--   name = "+Code coverage",
--   "<Cmd>lua coverage.load()<CR><Cmd>lua coverage.show()<CR>"
-- }

-- After changing plugin config exit and reopen LunarVim, Run :PackerInstall :PackerCompile
lvim.builtin.alpha.active = true
lvim.builtin.alpha.mode = "dashboard"
lvim.builtin.terminal.active = true
lvim.builtin.nvimtree.setup.view.side = "left"
lvim.builtin.nvimtree.setup.renderer.icons.show.git = false

-- if you don't want all the parsers change this to a table of the ones you want
lvim.builtin.treesitter.ensure_installed = {
  "bash",
  "c",
  "javascript",
  "json",
  "lua",
  "python",
  "typescript",
  "tsx",
  "css",
  "rust",
  "java",
  "yaml",
  "hcl",
  "go",
  "gomod",
  "gowork",
}

lvim.builtin.treesitter.ignore_install = { "haskell" }
lvim.builtin.treesitter.highlight.enable = true

-- generic LSP settings

-- -- make sure server will always be installed even if the server is in skipped_servers list
lvim.lsp.installer.setup.ensure_installed = {
  "lua_ls",
  "jsonls",
  "bashls",
  "gopls",
  "terraformls",
  "ansiblels",
  "pyright",
  "bufls",
  "groovyls",
}

-- -- change UI setting of `LspInstallInfo`
-- lvim.lsp.installer.setup.ui.check_outdated_servers_on_open = false
-- -- see <https://github.com/williamboman/nvim-lsp-installer#default-configuration>
-- lvim.lsp.installer.setup.ui.rounded = "rounded"
-- lvim.lsp.installer.setup.ui.keymaps = {
--   uninstall_server = "d",
--   toggle_server_expand = "o",
-- }

-- ---@usage disable automatic installation of servers
-- lvim.lsp.installer.setup.automatic_installation = false

-- ---configure a server manually. !!Requires `:LvimCacheReset` to take effect!!
-- ---see the full default list `:lua print(vim.inspect(lvim.lsp.automatic_configuration.skipped_servers))`
-- vim.list_extend(lvim.lsp.automatic_configuration.skipped_servers, { "pyright" })
-- local opts = {} -- check the lspconfig documentation for a list of all possible options
-- require("lvim.lsp.manager").setup("pyright", opts)

-- ---remove a server from the skipped list, e.g. eslint, or emmet_ls. !!Requires `:LvimCacheReset` to take effect!!
-- ---`:LvimInfo` lists which server(s) are skipped for the current filetype
-- lvim.lsp.automatic_configuration.skipped_servers = vim.tbl_filter(function(server)
--   return server ~= "emmet_ls"
-- end, lvim.lsp.automatic_configuration.skipped_servers)

-- -- you can set a custom on_attach function that will be used for all the language servers
-- -- See <https://github.com/neovim/nvim-lspconfig#keybindings-and-completion>
-- lvim.lsp.on_attach_callback = function(client, bufnr)
--   local function buf_set_option(...)
--     vim.api.nvim_buf_set_option(bufnr, ...)
--   end
--   --Enable completion triggered by <c-x><c-o>
--   buf_set_option("omnifunc", "v:lua.vim.lsp.omnifunc")
-- end

-- set a formatter, this will override the language server formatting capabilities (if it exists)
local formatters = require "lvim.lsp.null-ls.formatters"
formatters.setup {
  { command = "black", filetypes = { "python" } },
  { command = "isort", filetypes = { "python" } },
  {
    -- each formatter accepts a list of options identical to https://github.com/jose-elias-alvarez/null-ls.nvim/blob/main/doc/BUILTINS.md#Configuration
    command = "prettier",
    ---@usage arguments to pass to the formatter
    -- these cannot contain whitespaces, options such as `--line-width 80` become either `{'--line-width', '80'}` or `{'--line-width=80'}`
    extra_args = { "--print-with", "100" },
    ---@usage specify which filetypes to enable. By default a providers will attach to all the filetypes it supports.
    filetypes = { "typescript", "typescriptreact" },
  },
  { command = "gofumpt",       filetypes = { "go" }, },
  {
    command = "golines",
    extra_args = { "-m 100" },
    filetypes = { "go" },
  },
  { command = "terraform_fmt", filetypes = { "terraform", "tf", "hcl", "tfvars" }, },
}

-- set additional linters
local linters = require "lvim.lsp.null-ls.linters"
linters.setup {
  {
    command = "flake8",
    filetypes = { "python" },
  },
  -- Commented out because bashls uses shellcheck by default
  -- {
  --   -- each linter accepts a list of options identical to https://github.com/jose-elias-alvarez/null-ls.nvim/blob/main/doc/BUILTINS.md#Configuration
  --   command = "shellcheck",
  --   ---@usage arguments to pass to the formatter
  --   -- these cannot contain whitespaces, options such as `--line-width 80` become either `{'--line-width', '80'}` or `{'--line-width=80'}`
  --   extra_args = { "--severity", "warning" },
  -- },
  {
    command = "codespell",
    ---@usage specify which filetypes to enable. By default a providers will attach to all the filetypes it supports.
    -- filetypes = { "javascript", "python" },
  },
  {
    command = "golangci-lint",
    filetypes = { "go" },
  },
}

-- Additional Plugins
lvim.plugins = {
  {
    "folke/trouble.nvim",
    cmd = "TroubleToggle",
  },
  {
    "ThePrimeagen/refactoring.nvim",
    requires = {
      { "nvim-lua/plenary.nvim" },
      { "nvim-treesitter/nvim-treesitter" }
    }
  },
  {
    'echasnovski/mini.surround',
    branch = 'stable',
  },
  {
    "nvim-pack/nvim-spectre",
    requires = {
      "nvim-lua/plenary.nvim",
    }
  },
  {
    "andythigpen/nvim-coverage",
    requires = {
      "nvim-lua/plenary.nvim",
    }
  },
}

-- -- Configure refactoring for languages
require('refactoring').setup({
  -- prompt for return type
  prompt_func_return_type = {
    go = true,
  },
  -- prompt for function parameters
  prompt_func_param_type = {
    go = true,
  },
})

require('mini.surround').setup({
  -- Add custom surroundings to be used on top of builtin ones. For more
  -- information with examples, see `:h MiniSurround.config`.
  custom_surroundings = nil,
  -- Duration (in ms) of highlight when calling `MiniSurround.highlight()`
  highlight_duration = 500,
  -- Module mappings. Use `''` (empty string) to disable one.
  mappings = {
    add = 'sa',            -- Add surrounding in Normal and Visual modes
    delete = 'sd',         -- Delete surrounding
    find = 'sf',           -- Find surrounding (to the right)
    find_left = 'sF',      -- Find surrounding (to the left)
    highlight = 'sh',      -- Highlight surrounding
    replace = 'sr',        -- Replace surrounding
    update_n_lines = 'sn', -- Update `n_lines`
    suffix_last = 'l',     -- Suffix to search with "prev" method
    suffix_next = 'n',     -- Suffix to search with "next" method
  },
  -- Number of lines within which surrounding is searched
  n_lines = 20,
  -- How to search for surrounding (first inside current line, then inside
  -- neighborhood). One of 'cover', 'cover_or_next', 'cover_or_prev',
  -- 'cover_or_nearest', 'next', 'prev', 'nearest'. For more details,
  -- see `:h MiniSurround.config`.
  search_method = 'cover',
})
-- load refactoring telescope extension
require("telescope").load_extension("refactoring")
vim.cmd("set linebreak wrap")
-- Autocommands (https://neovim.io/doc/user/autocmd.html)
vim.api.nvim_create_autocmd("BufEnter", {
  pattern = { "*.json", "*.jsonc" },
  -- enable wrap mode for json files only
  command = "setlocal wrap",
})
vim.api.nvim_create_autocmd("FileType", {
  pattern = "zsh",
  callback = function()
    -- let treesitter use bash highlight for zsh files as well
    require("nvim-treesitter.highlight").attach(0, "bash")
  end,
})
-- Setup bufls with default values for protobuff
require("lspconfig").bufls.setup {
  cmd = { "bufls", "serve" },
  filetypes = { "proto" },
  root_pattern = { "buf.work.yaml", ".git" }
}

-- Groovy LSP setup
require("lspconfig").groovyls.setup {
  cmd = { "java", "-jar", "/usr/local/bin/groovy-language-server-all.jar" },
  filetypes = { "groovy", "Jenkinsfile" },
}

-- Setup test coverage display information
require("coverage").setup({
  commands = true, -- create commands
  highlights = {
    -- customize highlight groups created by the plugin
    covered = { fg = "#C3E88D" }, -- supports style, fg, bg, sp (see :h highlight-gui)
    uncovered = { fg = "#F07178" },
  },
  signs = {
    -- use your own highlight groups or text markers
    covered = { hl = "CoverageCovered", text = "▎" },
    uncovered = { hl = "CoverageUncovered", text = "▎" },
  },
  summary = {
    -- customize the summary pop-up
    min_coverage = 80.0, -- minimum coverage threshold (used for highlighting)
  },
  lang = {
    -- customize language specific settings
  },
})

-- Lua LSP setup
require("lspconfig").lua_ls.setup {
  settings = {
    Lua = {
      runtime = {
        -- Tell the language server which version of Lua you're using (most likely LuaJIT in the case of Neovim)
        version = 'LuaJIT',
      },
      diagnostics = {
        -- Get the language server to recognize the `vim` global
        globals = { 'vim' },
      },
      workspace = {
        -- Make the server aware of Neovim runtime files
        library = vim.api.nvim_get_runtime_file("", true),
      },
      -- Do not send telemetry data containing a randomized but unique identifier
      telemetry = {
        enable = false,
      },
    },
  },
}

-- Go LSP setup
local lspconfig = require "lspconfig"
local util = require "lspconfig/util"
lspconfig.gopls.setup {
  cmd = { "gopls", "serve" },
  filetypes = { "go", "gomod" },
  root_dir = util.root_pattern("go.work", "go.mod", ".git"),
  settings = {
    gopls = {
      analyses = {
        unusedparams = true,
      },
      staticcheck = true,
    },
  },
}
vim.api.nvim_create_autocmd('BufWritePre', {
  pattern = '*.go',
  callback = function()
    vim.lsp.buf.code_action({ context = { only = { 'source.organizeImports' } }, apply = true })
  end
})
