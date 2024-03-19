return {
	-- tools
	{
		"williamboman/mason.nvim",
		opts = function(_, opts)
			vim.list_extend(opts.ensure_installed, {
				"stylua",
				"shellcheck",
				"shfmt",
				"flake8",
				"luacheck",
				"tailwindcss-language-server",
				"css-lsp",
				"html-lsp",
				"htmx-lsp",
				"gopls",
				"golangci-lint",
				"gofumpt",
				"goimports",
				"zls",
				"codespell",
				"templ",
				"helm-ls",
				"bash-language-server",
			})
		end,
	},

	-- the opts function can also be used to change the default opts:
	{
		"nvim-lualine/lualine.nvim",
		event = "VeryLazy",
		opts = function(_, opts)
			table.insert(opts.sections.lualine_x, "😄")
		end,
	},

	-- or you can return new options to override all the defaults
	{
		"nvim-lualine/lualine.nvim",
		event = "VeryLazy",
		opts = function()
			return {
				--[[add your custom lualine config here]]
			}
		end,
	},
	-- lsp servers
	{
		"neovim/nvim-lspconfig",
		opts = {
			servers = {
				cssls = {
					filetypes = { "css", "scss", "less" },
				},
				htmx = {
					filetypes = { "html", "templ" },
				},
				tailwindcss = {
					root_dir = function(...)
						return require("lspconfig.util").root_pattern(".git")(...)
					end,
					filetypes = {
						"aspnetcorerazor",
						"astro",
						"astro-markdown",
						"blade",
						"clojure",
						"django-html",
						"htmldjango",
						"edge",
						"eelixir",
						"elixir",
						"ejs",
						"erb",
						"eruby",
						"gohtml",
						"gohtmltmpl",
						"haml",
						"handlebars",
						"hbs",
						"html",
						"html-eex",
						"heex",
						"jade",
						"leaf",
						"liquid",
						"markdown",
						"mdx",
						"mustache",
						"njk",
						"nunjucks",
						"php",
						"razor",
						"slim",
						"twig",
						"css",
						"less",
						"postcss",
						"sass",
						"scss",
						"stylus",
						"sugarss",
						"javascript",
						"javascriptreact",
						"reason",
						"rescript",
						"typescript",
						"typescriptreact",
						"vue",
						"svelte",
						"templ",
					},
				},
			},
		},
	},
}
