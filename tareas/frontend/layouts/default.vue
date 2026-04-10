<template>
  <div>
    <!-- Navbar -->
    <nav class="navbar is-primary" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a class="navbar-item has-text-weight-bold is-size-4" href="/">
          📋 TaskFlow
        </a>
      </div>
    </nav>

    <!-- Contenido principal -->
    <div class="container py-6">
      <slot />
    </div>
  </div>
</template>