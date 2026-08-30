return {
	{
		"williamboman/mason.nvim",
		opts = function(_, opts)
			opts.ensure_installed = opts.ensure_installed or {}
			vim.list_extend(opts.ensure_installed, { "expert" })
		end,
	},

	"neovim/nvim-lspconfig",
	opts = {
		servers = {
			-- Disable legacy ElixirLS if you use the LazyVim elixir extra
			elixirls = {
				enabled = false,
			},
			lexical = {
				enabled = false,
			},
			nextls = {
				enabled = false,
			},

			-- Enable Expert
			expert = {
				cmd = { "expert" },
				filetypes = { "elixir", "eelixir", "heex", "surface" },
				root_dir = function(fname)
					return require("lspconfig.util").root_pattern("mix.exs", ".git")(fname)
				end,
				settings = {},
			},
		},
		setup = {
			-- Custom setup handler in case mason-lspconfig hasn't mapped expert natively yet
			expert = function(_, opts)
				local configs = require("lspconfig.configs")
				if not configs.expert then
					configs.expert = {
						default_config = {
							cmd = opts.cmd,
							filetypes = opts.filetypes,
							root_dir = opts.root_dir,
							settings = opts.settings,
						},
					}
				end
				require("lspconfig").expert.setup(opts)
				return true
			end,
		},
	},
}
